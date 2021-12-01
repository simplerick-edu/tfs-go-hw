package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderEventsStorage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *OrderEventsStorage {
	return &OrderEventsStorage{pool}
}
