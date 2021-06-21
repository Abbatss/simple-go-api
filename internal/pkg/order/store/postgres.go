package store

import (
	"context"
	"errors"
	"github.com/Abbatss/TestGo/internal/pkg/order/order_errors"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (p *Postgres) Insert(ctx context.Context, order *Order) error {
	if order == nil {
		return errors.New("order entity is nil")
	}
	stmt := `
INSERT INTO
    orders (order_id, user_id)
VALUES($1, $2);
`
	cmd, err := p.pool.Exec(ctx, stmt, order.ID, order.UserID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return order_errors.ErrOrderNotFound
	}

	return nil
}

func (p *Postgres) GetByUser(ctx context.Context, userID string) ([]*Order, error) {
	stmt := `
SELECT
       order_id,
       user_id
FROM orders
WHERE user_id = $1;
`
	rows, err := p.pool.Query(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (p *Postgres) Get(ctx context.Context, orderID string) (*Order, error) {
	stmt := `
SELECT 
       order_id,
       user_id
FROM orders
WHERE order_id = $1;
`
	var order = &Order{}
	err := p.pool.QueryRow(ctx, stmt, orderID).Scan(&order.ID, &order.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, order_errors.ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}
