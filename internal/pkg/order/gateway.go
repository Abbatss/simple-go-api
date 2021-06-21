package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/Abbatss/TestGo/internal/pkg/order/store"
	"github.com/google/uuid"
	"strings"
	"time"
)

var timeout = time.Second * 3

type Gateway struct {
	storage Storage
}

type Storage interface {
	Get(ctx context.Context, orderID string) (*store.Order, error)
	GetByUser(ctx context.Context, userID string) ([]*store.Order, error)
	Insert(ctx context.Context, order *store.Order) error
}

func New(storage Storage) *Gateway {
	return &Gateway{storage: storage}
}

func (g *Gateway) GetOrder(ctx context.Context, entityID string) (*store.Order, error) {
	if IsNilOrEmpty(entityID) {
		return nil, errors.New("orderID is nil or empty")
	}
	cCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	order, err := g.storage.Get(cCtx, entityID)
	if err != nil {
		return nil, fmt.Errorf("can't get order from db for orderID:%s, error:%w", entityID, err)
	}

	return order, nil
}

func (g *Gateway) GetOrdersByUser(ctx context.Context, userID string) ([]*store.Order, error) {
	if IsNilOrEmpty(userID) {
		return nil, errors.New("userID is nil or empty")
	}
	cCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	orders, err := g.storage.GetByUser(cCtx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db for userID:%s, error:%w", userID, err)
	}

	return orders, nil
}

func (g *Gateway) InsertOrder(ctx context.Context, userID string) (*store.Order, error) {
	if IsNilOrEmpty(userID) {
		return nil, errors.New("userID is nil or empty")
	}
	cCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	order := &store.Order{
		ID:     uuid.NewString(),
		UserID: userID,
	}
	err := g.storage.Insert(cCtx, order)
	if err != nil {
		return nil, fmt.Errorf("can't insert order to db for userID:%s, error:%w", userID, err)
	}

	return order, nil
}

func IsNilOrEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
