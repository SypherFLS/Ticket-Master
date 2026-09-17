CREATE TABLE tickets (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'new',
    user_id BIGINT NOT NULL,

    CONSTRAINT tickets_status_check
        CHECK (status IN ('new', 'pending', 'close')),
    CONSTRAINT tickets_priority_check
        CHECK (priority IN ('high', 'middle', 'low')),
    CONSTRAINT fk_tickets_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);