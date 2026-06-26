-- +goose Up
CREATE TABLE IF NOT EXISTS notifications
(
    id          UUID        NOT NULL PRIMARY KEY,
    user_id     UUID        NOT NULL,
    type        TEXT        NOT NULL,
    payload     JSONB       NOT NULL,
    send_at     TIMESTAMP   NOT NULL DEFAULT NOW(),
    is_send     BOOL        NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
