package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type postgresStoreTestSuite struct {
	suite.Suite
	ctx context.Context

	pgxPool       *pgxpool.Pool
	postgresStore *Postgres
}

func clearDBRecords(suite *postgresStoreTestSuite) {
	// clear db data in some way. Otherwise another test can get this record on next review.
	_, err := suite.pgxPool.Exec(context.Background(), "TRUNCATE TABLE orders;")
	suite.Nil(err)
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(postgresStoreTestSuite))
}

func (suite *postgresStoreTestSuite) SetupTest() {
	require := suite.Require()

	suite.ctx = context.Background()

	suite.pgxPool = setup(require)
	store := NewPostgres(suite.pgxPool)
	suite.postgresStore = store
}

func setup(require *require.Assertions) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig("host=localhost port=5432 dbname=test user=user password=password")
	require.NoError(err)
	config.ConnConfig.Logger = zapadapter.NewLogger(zap.NewNop())

	conn, err := pgxpool.ConnectConfig(context.TODO(), config)
	require.NoError(err)

	return conn
}

func (suite *postgresStoreTestSuite) TestInsert_Get() {
	suite.Run(
		"should return notfound when no order in DB", func() {
			defer clearDBRecords(suite)

			order, err := suite.postgresStore.Get(suite.ctx, uuid.NewString())
			suite.Error(err)
			suite.Nil(order)
		},
	)

	suite.Run(
		"should return error when order is nil", func() {
			defer clearDBRecords(suite)

			err := suite.postgresStore.Insert(suite.ctx, nil)
			suite.Error(err)
		},
	)
	suite.Run(
		"should insert and get order successfully", func() {
			defer clearDBRecords(suite)
			order := &Order{
				ID:     "1",
				UserID: "2",
			}

			err := suite.postgresStore.Insert(suite.ctx, order)
			suite.Nil(err)

			dbOrder, err := suite.postgresStore.Get(suite.ctx, order.ID)

			suite.Nil(err)
			suite.EqualValues(order, dbOrder)

		},
	)

	suite.Run(
		"should fail if insert same order 2 times", func() {
			defer clearDBRecords(suite)
			order := &Order{
				ID:     "1",
				UserID: "2",
			}

			err := suite.postgresStore.Insert(suite.ctx, order)
			suite.Nil(err)

			err = suite.postgresStore.Insert(suite.ctx, order)
			suite.Error(err)

		},
	)

}

func (suite *postgresStoreTestSuite) Test_GetByUser() {
	suite.Run(
		"should return empty list if no user orders", func() {
			defer clearDBRecords(suite)

			orders, err := suite.postgresStore.GetByUser(suite.ctx, uuid.NewString())
			suite.Nil(err)
			suite.Len(orders, 0)
		},
	)

	suite.Run(
		"should return only user orders", func() {
			defer clearDBRecords(suite)
			order1 := &Order{ID: "1", UserID: "2"}
			order2 := &Order{ID: "2", UserID: "2"}
			order3 := &Order{ID: "3", UserID: "3"}

			err := suite.postgresStore.Insert(suite.ctx, order1)
			suite.Nil(err)
			err = suite.postgresStore.Insert(suite.ctx, order2)
			suite.Nil(err)
			err = suite.postgresStore.Insert(suite.ctx, order3)
			suite.Nil(err)

			orders, err := suite.postgresStore.GetByUser(suite.ctx, "2")
			suite.Nil(err)
			suite.Len(orders, 2)
			suite.EqualValues(order1, orders[0])
			suite.EqualValues(order2, orders[1])
		},
	)

}
