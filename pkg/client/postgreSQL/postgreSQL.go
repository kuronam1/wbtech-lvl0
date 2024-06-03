package postgreSQL

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"wbLvL0/internal/config"
)

const (
	connectPGTimeout = 5 * time.Second
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
}

func NewClient(cfg config.PG) (*pgxpool.Pool, error) {
	ctx, cancelPG := context.WithTimeout(context.Background(), connectPGTimeout)
	defer cancelPG()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Login,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SslMode)

	c, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	c.MaxConns = 20 //int32(runtime.NumCPU())
	c.MinConns = 5

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
