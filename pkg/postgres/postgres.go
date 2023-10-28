package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mxmrykov/L0/config"
)

func Connect(cfg *config.PG) (*pgxpool.Pool, error) {
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	return pgxpool.Connect(context.Background(), connection)
}
