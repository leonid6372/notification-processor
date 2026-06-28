-- +goose Up
CREATE TYPE notification_type AS ENUM ('order_created', 'payment_received', 'order_shipped');
CREATE TYPE notification_status AS ENUM ('created', 'in_progress', 'completed', 'failed');

CREATE TABLE IF NOT EXISTS notifications
(
    id          UUID                    NOT NULL PRIMARY KEY,
    user_id     UUID                    NOT NULL,
    type        notification_type       NOT NULL,
    payload     JSONB                   NOT NULL,
    send_at     TIMESTAMP               NOT NULL DEFAULT NOW(),
    status      notification_status     NOT NULL DEFAULT 'created',
    tries_count INT                     NOT NULL DEFAULT 0,
    started_at  TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_status_send_at ON notifications (status, send_at);
CREATE INDEX IF NOT EXISTS idx_notifications_status_started_at ON notifications (status, started_at);

-- +goose Down
DROP INDEX IF EXISTS idx_notifications_status_send_at;
DROP INDEX IF EXISTS idx_notifications_status_started_at;
DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS notification_status;
DROP TYPE IF EXISTS notification_type;
