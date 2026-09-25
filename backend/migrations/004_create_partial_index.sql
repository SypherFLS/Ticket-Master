-- +goose Up

CREATE UNIQUE INDEX idx_one_pending_ticket_per_operator
ON tickets (operator_id)
WHERE status = 'pending' AND operator_id IS NOT NULL;

-- +goose Down

DROP INDEX idx_one_pending_ticket_per_operator;