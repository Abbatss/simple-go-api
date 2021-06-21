package pgx

import (
	"context"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

var _ svc.Worker = (*Worker)(nil)

type Worker struct {
	pool *pgxpool.Pool
}

func (w *Worker) Init(logger *zap.Logger) error {
	return nil
}

func (w *Worker) Run() error {
	return nil
}

func (w *Worker) Terminate() error {
	w.pool.Close()
	return nil
}

// connString like "host=localhost port=5432 dbname=test user=user password=password"
func Connect(connString string) (*Worker, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.ConnConfig.Logger = zapadapter.NewLogger(zap.NewNop())

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &Worker{pool: conn}, nil
}
func (w *Worker) Pool() *pgxpool.Pool {
	return w.pool
}
