package repository

import (
	"context"
	"fmt"
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
        CREATE TABLE IF NOT EXISTS main (
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
  			sm_id INTEGER,
            date_created VARCHAR(255),
            oof_shard VARCHAR(255)
        )
    `)

	return err
}

func (r *Repo) SaveOrder(order models.Order) error {
	_, err := r.pool.Exec(context.Background(), `
        INSERT INTO main (order_uid, track_number, entry, delivery_info, payment_info, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `, order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payment, order.Items, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard)

	return err
}

func (r *Repo) GetALl() ([]models.Order, error) {

	res, err := r.pool.Query(context.Background(),
		`SELECT * FROM main`)

	if err != nil {
		fmt.Printf("Error at getting all orders: %v\n", err)
		return nil, err
	}

	defer res.Close()

	var orders []models.Order
	for res.Next() {
		var execOrder models.Order
		err := res.Scan(
			&execOrder.OrderUid,
			&execOrder.TrackNumber,
			&execOrder.Entry,
			&execOrder.Delivery,
			&execOrder.Payment,
			&execOrder.Items,
			&execOrder.Locale,
			&execOrder.InternalSignature,
			&execOrder.CustomerId,
			&execOrder.DeliveryService,
			&execOrder.ShardKey,
			&execOrder.SmId,
			&execOrder.DateCreated,
			&execOrder.OofShard,
		)

		if err != nil {
			fmt.Printf("Error at parsing preloading: %v\n", err)
		}

		orders = append(orders, execOrder)
	}

	if err := res.Err(); err != nil {
		fmt.Printf("Error at final getting orders: %v\n", err)
		return nil, err
	}
	fmt.Println("Orders preloaded DB -> Cache")
	return orders, nil
}

// useless, т.к. все заказы разом подгружаются в кэш и данные http запросов читаются из кэша
func (r *Repo) GetOrder(uid string) (models.Order, error) {

	var execOrder models.Order

	err := r.pool.QueryRow(context.Background(), `SELECT * FROM main WHERE order_uid = $1`, uid).
		Scan(
			&execOrder.OrderUid,
			&execOrder.TrackNumber,
			&execOrder.Entry,
			&execOrder.Delivery,
			&execOrder.Payment,
			&execOrder.Items,
			&execOrder.Locale,
			&execOrder.InternalSignature,
			&execOrder.CustomerId,
			&execOrder.DeliveryService,
			&execOrder.ShardKey,
			&execOrder.SmId,
			&execOrder.DateCreated,
			&execOrder.OofShard,
		)

	if err != nil {
		return execOrder, err
	}

	return execOrder, nil

}
