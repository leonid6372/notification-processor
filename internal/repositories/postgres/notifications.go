package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/pkg/errs"
)

type notificationsRepo struct {
	psql *pgxpool.Pool
}

func NewNotificationsRepo(pool *pgxpool.Pool) domains.NotificationsRepo {
	return &notificationsRepo{
		psql: pool,
	}
}

func (n *notificationsRepo) CreateNotification(ctx context.Context, notification *domains.Notification) error {
	query := `INSERT INTO notifications (id, user_id, type, payload, send_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := n.psql.Exec(
		ctx,
		query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.RawPayload,
		notification.SendAt,
	)
	if err != nil {
		return errs.NewStack(err)
	}

	return nil
}

func (n *notificationsRepo) GetNotificationsToSend(ctx context.Context, limit int) ([]*domains.Notification, error) {
	query := `UPDATE notifications
			SET status = 'in_progress', tries_count = tries_count + 1, started_at = NOW()
			WHERE id IN (
						SELECT id
						FROM notifications
						WHERE status = 'created' AND send_at <= NOW()
						ORDER BY send_at ASC
						LIMIT $1
					)
			RETURNING *`

	rows, err := n.psql.Query(ctx, query, limit)
	if err != nil {
		return nil, errs.NewStack(err)
	}
	defer rows.Close()

	notifications := make([]*domains.Notification, 0)
	for rows.Next() {
		notification := new(Notification)

		err := rows.Scan(notification)
		if err != nil {
			return nil, errs.NewStack(err)
		}

		domainNotification, err := notification.ToDomain()
		if err != nil {
			return nil, errs.NewStack(err)
		}

		notifications = append(notifications, domainNotification)
	}

	return notifications, nil
}

func (n *notificationsRepo) UpdateNotificationStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE notifications SET status = $1 WHERE id = $2`

	if _, err := n.psql.Exec(ctx, query, status, id); err != nil {
		return errs.NewStack(err)
	}

	return nil
}

func (n *notificationsRepo) ResetZombieNotifications(ctx context.Context, timeout time.Duration) error {
	query := `UPDATE notifications
			SET status = 'created'
			WHERE status = 'in_progress' AND started_at <= NOW() - $1::interval`

	if _, err := n.psql.Exec(ctx, query, timeout.String()); err != nil {
		return errs.NewStack(err)
	}

	return nil
}
