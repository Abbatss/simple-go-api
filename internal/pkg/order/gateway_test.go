package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/Abbatss/TestGo/internal/pkg/order/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var _ Storage = (*storageMock)(nil)

type storageMock struct {
	GetFunc       func(ctx context.Context, orderID string) (*store.Order, error)
	GetByUserFunc func(ctx context.Context, userID string) ([]*store.Order, error)
	InsertFunc    func(ctx context.Context, order *store.Order) error
}

func (s *storageMock) Get(ctx context.Context, orderID string) (*store.Order, error) {
	return s.GetFunc(ctx, orderID)
}

func (s *storageMock) GetByUser(ctx context.Context, userID string) ([]*store.Order, error) {
	return s.GetByUserFunc(ctx, userID)
}

func (s *storageMock) Insert(ctx context.Context, order *store.Order) error {
	return s.InsertFunc(ctx, order)
}

func TestGateway_GetOrder(t *testing.T) {
	orderID := uuid.NewString()
	tests := []struct {
		name    string
		want    *store.Order
		store   *storageMock
		wantErr error
		orderID string
	}{
		{
			name:    "empty orderID should return error",
			wantErr: errors.New("orderID is nil or empty"),
		},
		{
			name: "should return error if storage return error",
			store: &storageMock{GetFunc: func(ctx context.Context, orderID string) (*store.Order, error) {
				return nil, errors.New("error")
			}},
			orderID: orderID,

			wantErr: fmt.Errorf("can't get order from db for orderID:%s, error:%w", orderID, errors.New("error")),
		},
		{
			name: "should return order from store",
			store: &storageMock{GetFunc: func(ctx context.Context, orderID string) (*store.Order, error) {
				return &store.Order{
					ID:     "ID",
					UserID: "UserID",
				}, nil
			}},
			want: &store.Order{
				ID:     "ID",
				UserID: "UserID",
			},
			orderID: orderID,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := &Gateway{
				storage: tt.store,
			}
			got, err := g.GetOrder(context.TODO(), tt.orderID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestGateway_GetOrdersByUser(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*store.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gateway{
				storage: tt.fields.storage,
			}
			got, err := g.GetOrdersByUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrdersByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrdersByUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGateway_InsertOrder(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *store.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gateway{
				storage: tt.fields.storage,
			}
			got, err := g.InsertOrder(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNilOrEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", " ", true},
		{"non empty string", " s ", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNilOrEmpty(tt.s); got != tt.want {
				t.Errorf("IsNilOrEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	storage := &storageMock{}
	gw := New(storage)
	assert.NotNil(t, gw)
	assert.Equal(t, storage, gw.storage)
}
