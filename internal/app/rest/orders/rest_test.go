package orders

import (
	"context"
	"github.com/Abbatss/TestGo/internal/pkg/order/store"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
)

const (
	zoneName = "stockholm"
	roleType = "MDS"
)

type orderGatewayMock struct {
	GetOrderFunc func(ctx context.Context, entityID string) (*store.Order, error)
}

func (o *orderGatewayMock) GetOrder(ctx context.Context, entityID string) (*store.Order, error) {
	return o.GetOrderFunc(ctx, entityID)
}

type ordersControllerTestSuite struct {
	suite.Suite
	recorder      *httptest.ResponseRecorder
	ordersGateway *orderGatewayMock
	router        chi.Router
}

func TestPartnersController(t *testing.T) {
	suite.Run(t, new(ordersControllerTestSuite))
}

func (suite *ordersControllerTestSuite) SetupTest() {
	suite.ordersGateway = &orderGatewayMock{}

	suite.recorder = httptest.NewRecorder()

	controller := New(zap.NewNop(), suite.ordersGateway)
	suite.router = chi.NewRouter()
	subRouter := chi.NewRouter()
	subRouter.Route("/", func(publicGroup chi.Router) {
		controller.RegisterRoutes(publicGroup)
	})
	suite.router.Mount("/", subRouter)
}

func (suite *ordersControllerTestSuite) TestGetOrder() {
	require := suite.Require()
	suite.ordersGateway.GetOrderFunc = func(ctx context.Context, entityID string) (*store.Order, error) {
		suite.Require().Equal("1", entityID)
		return &store.Order{
			ID:     "ID",
			UserID: "UserID",
		}, nil
	}
	request := httptest.NewRequest(
		http.MethodGet,
		"/orders/1",
		nil,
	)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()

	require.Equal(http.StatusOK, response.StatusCode)
}
