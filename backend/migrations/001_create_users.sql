-- +goose Up

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
        CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    role VARCHAR(20) NOT NULL DEFAULT 'user',

    CONSTRAINT users_role_check
        CHECK (role IN ('user', 'operator', 'admin'))
);

-- +goose Down

DROP TABLE users;