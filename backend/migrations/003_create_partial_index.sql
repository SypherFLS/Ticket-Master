-- +goose Up

CREATE INDEX idx_tickets_new_queue
ON tickets (created_at, id)
WHERE status = 'new';

-- +goose Down 

DROP INDEX idx_tickets_new_queue;