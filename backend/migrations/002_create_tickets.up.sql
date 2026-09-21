CREATE TABLE tickets (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    operator_id BIGINT REFERENCES users(id),
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'new',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    claimed_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,

    CONSTRAINT tickets_status_check
        CHECK (status IN ('new', 'pending', 'closed')),
    CONSTRAINT fk_tickets_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);

