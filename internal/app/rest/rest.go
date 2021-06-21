package rest

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

type SubController interface {
	RegisterRoutes(router chi.Router)
}

type Controller struct {
	subcontrollers []SubController
}

// New creates a new REST Controller for the service.
func New(subcontrollers ...SubController) *Controller {
	return &Controller{
		subcontrollers: subcontrollers,
	}
}

// Controller interface --------------------------------------------------------

func (c *Controller) Init(_ *zap.Logger) error {
	return nil
}

func (c *Controller) SetupRouter(router chi.Router) error {
	router.Route("/", func(r chi.Router) {
		r.Use(c.RequireLogin)
		c.setupSubRouters(r)
	})
	return nil
}

func (c *Controller) Terminate() error {
	return nil
}
func (c *Controller) setupSubRouters(r chi.Router) {
	for _, subController := range c.subcontrollers {
		subController.RegisterRoutes(r)
	}
}

func (c *Controller) RequireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !c.isAuthenticated(r) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (c *Controller) isAuthenticated(r *http.Request) bool {
	return true
}
