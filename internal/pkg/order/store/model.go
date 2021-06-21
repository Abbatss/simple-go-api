package store

type Order struct {
	ID     string `db:"order_id"`
	UserID string `db:"user_id"`
}
