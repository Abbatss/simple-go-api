package orders

import (
	"context"
	"github.com/Abbatss/TestGo/internal/app/rest"
	"github.com/Abbatss/TestGo/internal/pkg/order/order_errors"
	"github.com/Abbatss/TestGo/internal/pkg/order/store"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type gateway interface {
	GetOrder(ctx context.Context, entityID string) (*store.Order, error)
}

var _ rest.SubController = (*Controller)(nil)

type Controller struct {
	log     *zap.Logger
	service gateway
}

// New creates a new REST Controller for the service.
func New(logger *zap.Logger, gateway gateway) *Controller {
	return &Controller{
		log:     logger,
		service: gateway,
	}
}

// Controller interface --------------------------------------------------------

func (c *Controller) Init(log *zap.Logger) error {
	c.log = log
	return nil
}

func (c *Controller) RegisterRoutes(router chi.Router) {
	router.Get("/orders/{id}", c.GetEntityByID)
}

func (c *Controller) Terminate() error {
	return nil
}

func (c *Controller) GetEntityByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	l := c.log.With(zap.String("id", id))

	res, err := c.service.GetOrder(r.Context(), id)
	if err == order_errors.ErrOrderNotFound {
		l.Debug("error when getting order", zap.Error(err))
		rest.JSON(w, http.StatusNotFound, nil)
		return
	}

	if err != nil {
		l.Error("error when getting order", zap.Error(err))
		rest.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	rest.JSON(w, http.StatusOK, mapToModel(res))
}

func mapToModel(res *store.Order) *Order {
	return &Order{
		ID:     res.ID,
		UserID: res.UserID,
	}
}
