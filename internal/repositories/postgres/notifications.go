package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonid6372/notification-processor/internal/domain"
)

type notificationsRepository struct {
	psql *pgxpool.Pool
}

func NewNotificationsRepository(pool *pgxpool.Pool) domain.NotificationsRepository {
	return &notificationsRepository{
		psql: pool,
	}
}
