package main

import (
	"fmt"
	"github.com/Abbatss/TestGo/internal/app/rest"
	"github.com/Abbatss/TestGo/internal/app/rest/orders"
	"github.com/Abbatss/TestGo/internal/pkg/order"
	"github.com/Abbatss/TestGo/internal/pkg/order/store"
	"github.com/Abbatss/TestGo/internal/pkg/workers/chiworker"
	"github.com/Abbatss/TestGo/internal/pkg/workers/pgx"
	"go.uber.org/zap"
	"time"

	"github.com/voi-oss/svc"
)

const (
	serviceName    = "test"
	serviceVersion = "snapshot"
)

// Environment overridable configs
type config struct {
	Env string `env:"APP_ENV envDefault:"dev""`
}

func main() {
	cfg := config{}

	// Read up global configs
	if err := svc.LoadFromEnv(&cfg); err != nil {
		panic(fmt.Sprintf("could not load configuration: %s", err))
	}

	// SVC supervisor Init
	options := []svc.Option{
		svc.WithTerminationGracePeriod(55 * time.Second),
		svc.WithTerminationWaitPeriod(30 * time.Second),
		svc.WithMetrics(),
		svc.WithHealthz(),
		svc.WithMetricsHandler(),
		svc.WithLogLevelHandlers(),
		svc.WithHTTPServer("9090"),
	}

	options = append(options, svc.WithStackdriverLogger(zap.DebugLevel))

	// SVC supervisor Init
	service, err := svc.New(serviceName, serviceVersion, options...)
	svc.MustInit(service, err)

	pgWorker, err := pgx.Connect("host=localhost port=5432 dbname=test user=user password=password")
	svc.MustInit(service, err)
	service.AddWorker("postgres", pgWorker)

	orderGW := order.New(store.NewPostgres(pgWorker.Pool()))
	ordersController := orders.New(service.Logger(), orderGW)
	restController := rest.New(ordersController)
	restWorker := chiworker.New(restController)

	service.AddWorker("rest", restWorker)

	service.Run()

}
