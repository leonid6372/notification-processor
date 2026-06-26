-- +goose Up
CREATE TYPE notification_type AS ENUM ('order_created', 'payment_received', 'order_shipped');

CREATE TABLE IF NOT EXISTS notifications
(
    id          UUID                    NOT NULL PRIMARY KEY,
    user_id     UUID                    NOT NULL,
    type        notification_type       NOT NULL,
    payload     JSONB                   NOT NULL,
    send_at     TIMESTAMP               NOT NULL DEFAULT NOW(),
    is_send     BOOL                    NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS notification_type;
