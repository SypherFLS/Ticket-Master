# Ticket-Master

Проект представляет собой backend-сервис для управления тикетами в формате мини-CRM/системы поддержки. Основная идея — разделить пользователей на роли и дать каждой роли только тот функционал, который ей нужен: пользователь создаёт тикеты, оператор обрабатывает очередь, администратор управляет ролями и имеет доступ к более широким операциям.

## Основная логика

Система строится вокруг трёх ключевых сущностей:
- User — пользователь системы
- Ticket — тикет, который создаётся пользователем и затем обрабатывается оператором
- Role — роль пользователя: user, operator, admin

Тикет проходит жизненный цикл:
- new — новый тикет, ожидает обработки
- pending — тикет взят в работу оператором
- closed — тикет закрыт и помечен как решённый

## Разделение ролей

- User:
  - может регистрироваться и входить в систему
  - может создавать тикеты
  - может смотреть только свои тикеты

- Operator:
  - может смотреть очередь новых тикетов
  - может взять следующий тикет в работу
  - может закрывать тикет после решения
  - может видеть свои назначенные тикеты

- Admin:
  - имеет расширенный доступ к тикетам
  - может управлять ролями пользователей
  - может выполнять действия, связанные с общей помощью и администрированием

Вся логика прав реализована через permission map. Это позволяет не завязываться на жёстком условии `if role == "admin"`, а хранить права как набор допустимых операций. Для каждого пользователя роли сопоставляется набор разрешений, а сервисный слой проверяет наличие нужного permission перед выполнением операции.

Пример структуры прав:

```go
var RolePermission = map[constants.UserRole]map[Permission]struct{}{
    constants.RoleUser: {
        PermissionTicketCreate:  {},
        PermissionTicketReadOwn: {},
    },
    constants.RoleOperator: {
        PermissionTicketOwnClaimed: {},
        PermissionTicketAssign:     {},
        PermissionTicketReadQueue:  {},
        PermissionTicketUpdate:     {},
    },
    constants.RoleAdmin: {
        PermissionTicketReadAll:   {},
        PermissionTicketReadQueue: {},
        PermissionUserManage:      {},
    },
}
```

Это типичный pattern "role-based access control" с упрощённым permission map, где доступ определяется по роли и набору разрешений.

## Архитектура проекта

Проект выполнен в монолитной архитектуре, но логически разделён на слои. Главная идея — не смешивать API-логику, бизнес-логику и доступ к данным.

### 1. API слой

API слой реализован на `net/http` без фреймворка. Есть маршрутизатор, middleware и handlers.

Основной вход в приложение:

```go
router := router.NewRouter(handler, JWTManager, cfg)
http.ListenAndServe(cfg.Server.Port, router)
```

В `router.go` создаются два маршрута:
- public — регистрация и логин
- private — все защищённые endpoint'ы, требующие JWT

```go
public.Handle("POST /register", http.HandlerFunc(h.RegisterHandler))
public.Handle("POST /login", http.HandlerFunc(h.LoginHandler))

private.Handle("POST /create", http.HandlerFunc(h.CreateTicketHandler))
private.Handle("GET /tickets", http.HandlerFunc(h.GetOwnTicketsHandler))
private.Handle("POST /claime", http.HandlerFunc(h.ClaimNextTicketHandler))
private.Handle("PATCH /close", http.HandlerFunc(h.CloseTicketHandler))
private.Handle("GET /tickets_queue", http.HandlerFunc(h.GetNewTicketsHandler))
```

### 2. Middleware слой

В проекте используется цепочка middleware, которая добавляет общий функционал для всех запросов:
- логирование
- восстановление после падения
- таймаут на запрос
- авторизация JWT

Это хороший пример паттерна "chain of responsibility". Каждый middleware оборачивает `http.Handler` и добавляет отдельную обязанность.

```go
privateChain := middlewares.CommonChain(
    middlewares.AuthMiddleware(jwtManager)(private),
    cfg.Server.Timeout,
)
```

`AuthMiddleware` принимает `Authorization: Bearer <token>`, валидирует JWT и кладёт `user_id` и `user_role` в контекст запроса. После этого последующие обработчики работают уже с авторизованным пользователем.

### 3. Сервисный слой

Сервисный слой отвечает за бизнес-логику. Он проверяет права пользователя, маппит DTO в модели и вызывает репозиторий.

Пример:

```go
func (s *Service) CreateTicketService(ctx context.Context, ticketDTO dto.TicketDTO, user_id int, user_Role constants.UserRole) (int, error) {
    if !auth.HasPermission(user_Role, auth.PermissionTicketCreate) {
        return 0, apperrors.NoPermission
    }

    ticket := dto.TicketToModel(ticketDTO, user_id)
    return ticket.ID, s.repo.CreateTicketRepo(ctx, &ticket)
}
```

Таким образом выделяется слой бизнес-правил: handler не должен сам решать, можно ли пользователю создавать тикет, а репозиторий не должен проверять роль пользователя.

### 4. Репозиторий

Репозиторий инкапсулирует работу с БД. На него уходит общий интерфейс:

```go
type Repository interface {
    RegisterRepo(ctx context.Context, user models.User) error
    GetLoginDataRepo(ctx context.Context, email string) (models.LoginResult, error)
    SetRoleRepo(ctx context.Context, email string, role constants.UserRole) error

    GetOwnTicketsRepo(ctx context.Context, user_id int, pagData dto.PaginationData) ([]models.Ticket, error)
    GetNewTickets(ctx context.Context, limit int) ([]models.Ticket, error)
    ClaimNextTicketRepo(ctx context.Context, operatorID int) (*models.Ticket, error)
    GetClaimedTicket(ctx context.Context, operator_id int) (models.Ticket, error)
    CreateTicketRepo(ctx context.Context, ticket *models.Ticket) error
    CloseTicketRepo(ctx context.Context, ticket_id int, operator_id int) error
}
```

Это пример паттерна Repository, который позволяет скрыть детали СУБД и подключить любую реализацию без изменения сервисного слоя.

## Механика работы запросов

### Регистрация и логин

При регистрации пользователь создаётся через `RegisterRepo`, а пароль хэшируется до записи в базу. После логина система генерирует JWT с `user_id` и `user_role`.

```go
func (j *JWTManager) Generate(userID int, UserRole string) (string, error) {
    claims := claims{
        UserID:   userID,
        UserRole: UserRole,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secret)
}
```

JWT использует HS256 и валидируется на каждом защищённом запросе через `AuthMiddleware`.

### Создание тикета

Пользователь отправляет запрос на `/api/create`. В обработчике берётся `user_id` из контекста, сервис проверяет permission `ticket:create`, затем конвертирует DTO в модель тикета и сохраняет в БД.

Логика очень простая и соответствует схеме:
- handler -> service -> repository -> database

### Просмотр своих тикетов

Пользователь может вызвать `GET /api/tickets`. Сервис сначала проверяет `PermissionTicketReadOwn`, затем запрашивает тикеты через `GetOwnTicketsRepo` с фильтрацией по `user_id` и пагинацией.

### Очередь тикетов для оператора

У оператора есть отдельный способ получения тикетов из очереди: `GetNewTickets`. Это запрос к базе по статусу `new`, отсортированный по времени создания.

### Взятие тикета в работу

Самое важное место в проекте — логика назначения тикета оператору.

```go
err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    result := tx.Where("status = ?", constants.StatusNew).
        Order("created_at ASC, id ASC").
        Clauses(clause.Locking{
            Strength: "UPDATE",
            Options:  "SKIP LOCKED",
        }).
        First(&ticket)
```

Здесь реализован паттерн конкурентного назначения тикетов:
- открывается транзакция
- выбирается первый свободный `new` тикет
- строка блокируется `FOR UPDATE SKIP LOCKED`
- тикет переводится в `pending`
- оператор записывается в поле `operator_id`
- фиксируется `claimed_at`

Это защищает систему от race condition, когда несколько операторов захотят взять один и тот же тикет одновременно.

### Закрытие тикета

Когда оператор решает вопрос, вызывается `CloseTicketService`. Проверяется permission `ticket:update`, затем в репозитории выполняется UPDATE по `id` и `operator_id`.

```go
res := r.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ? AND operator_id = ?", ticket_id, operator_id).Updates(map[string]any{
    "status":      constants.StatusClosed,
    "resolved_at": now,
})
```

То есть закрывать тикет может только тот оператор, который его взял в работу.

## Модели данных

### User

```go
type User struct {
    ID           int    `gorm:"primaryKey"`
    Name         string `gorm:"not null;unique"`
    PasswordHash string
    Email        string             `gorm:"not null;unique"`
    Role         constants.UserRole `gorm:"not null;default:user"`
    Tickets      []Ticket
}
```

### Ticket

```go
type Ticket struct {
    ID          int    `gorm:"primaryKey"`
    Title       string `gorm:"not null"`
    OperatorID  *int
    Description string
    Status      constants.TicketStatus `gorm:"not null,default:new"`
    CreatedAt   time.Time
    ClaimedAt   *time.Time
    ResolvedAt  *time.Time
    UserID      int  `gorm:"not null;index"`
    User        User `gorm:"foreignKey:UserID"`
}
```

Связь между пользователем и тикетом — один ко многим: один пользователь может иметь много тикетов, а тикет всегда привязан к создателю.

## База данных и миграции

В качестве СУБД используется PostgreSQL. Подключение и настройка конфигурируются через `config.yaml` и `.env`/переменные окружения.

Используется GORM для работы с БД, а миграции лежат в папке `backend/migrations`:
- `001_create_users.up.sql`
- `002_create_tickets.up.sql`
- `003_create_partial_index.up.sql`
- `004_create_partial_index.up.sql`

Это даёт возможность постепенно эволюционировать схему БД и поддерживать код и структуру данных в согласованном состоянии.

## Конфигурация и запуск

Проект умеет поднимать разные окружения через переменные `ENV`:
- `local`
- `dev`

В `main.go` выбирается нужный конфиг и инициализируется БД:

```go
config_env := ""
switch os.Getenv("ENV") {
case "local":
    config_env = "CONFIG_PATH_BACKEND"
case "dev":
    config_env = "CONFIG_PATH_DOCKER"
default:
    log.Fatal("wrong work env")
}
```

После инициализации создаются:
- DB connection
- JWT manager
- repository
- service
- handlers
- router

## Паттерны и подходы, которые используются в проекте

1. Repository pattern
   - всё взаимодействие с БД спрятано в `repository`
   - сервисы ничего не знают про SQL или GORM

2. Dependency injection
   - `repo` и `JWTManager` создаются в `main.go` и передаются по слоям

3. Middleware chain
   - логирование, таймаут, авторизация и восстановление ошибок объединены в цепочку

4. Role-based access control
   - права хранятся в `RolePermission`
   - проверка идёт в сервисе до бизнес-операции

5. Transaction + locking
   - для конкурирующего назначения тикетов используется транзакция и блокировка строки `SKIP LOCKED`

6. Context-based request flow
   - `context.Context` используется для передачи информации о пользователе и таймаутов запроса

7. Monolith with clear boundaries
   - есть логическое разделение на API, middleware, auth, repository, service, models, constants

## Краткое резюме

Ticket-Master — это небольшой, но уже довольно логично построенный монолитный сервис для поддержки тикетов. Он сочетает в себе:
- `net/http` API без тяжёлого фреймворка
- JWT авторизацию
- роль-ориентированную систему прав
- PostgreSQL + GORM
- миграции и конфиг
- middleware-цепочку для защиты и логирования
- конкурентную обработку очереди тикетов через транзакции и lock

Если смотреть на проект как на учебный/прототип, то он показывает хороший набор практик для backend-сервиса: чистые слои, авторизация, бизнес-логика, репозиторий, ограничения прав и обработка конкурентных сценариев.

```go
type Repository interface {
    RegisterRepo(ctx context.Context, user models.User) error
    GetLoginDataRepo(ctx context.Context, email string) (models.LoginResult, error)
    SetRoleRepo(ctx context.Context, email string, role constants.UserRole) error

    GetOwnTicketsRepo(ctx context.Context, user_id int, pagData dto.PaginationData) ([]models.Ticket, error)
    GetNewTickets(ctx context.Context, limit int) ([]models.Ticket, error)
    ClaimNextTicketRepo(ctx context.Context, operatorID int) (*models.Ticket, error)
    CreateTicketRepo(ctx context.Context, ticket *models.Ticket) error
    CloseTicketRepo(ctx context.Context, ticket_id uint, operator_id int) error
}
```

Это базовая схема, вокруг которой и строится вся система: пользователь создаёт тикет, оператор его забирает, администратор управляет ролями, а поведение ограничено через permissions и JWT.

