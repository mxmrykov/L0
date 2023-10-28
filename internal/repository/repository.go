package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mxmrykov/L0/internal/models"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repo {
	return &Repo{
		pool: pool,
	}
}

func (r *Repo) CreateTable() error {
	_, err := r.pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS orders (
            order_uid VARCHAR(255) PRIMARY KEY,
            track_number VARCHAR(255),
            entry VARCHAR(255),
            delivery_info JSONB,
            payment_info JSONB,
            items JSONB,
            locale VARCHAR(255),
            internal_signature VARCHAR(255),
  			customer_id VARCHAR(255),
  			delivery_service VARCHAR(255),
  			shardkey VARCHAR(255),
  			sm_id VARCHAR(255),
            date_created TIMESTAMP,
            oof_shard VARCHAR(255)
        )
    `)

	return err
}

func (r *Repo) SaveOrder(order models.Order) error {
	_, err := r.pool.Exec(context.Background(), `
        INSERT INTO orders (order_uid, track_number, entry, delivery_info, payment_info, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `, order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payment, order.Items, order.Locale, order.DateCreated, order.OofShard)

	return err
}
