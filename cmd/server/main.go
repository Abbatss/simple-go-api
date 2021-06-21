package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
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
	service.Run()

}

func setupPostgres(service *svc.SVC) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig("host=localhost port=5432 dbname=test user=user password=password")
	svc.MustInit(service, err)
	config.ConnConfig.Logger = zapadapter.NewLogger(zap.NewNop())

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	svc.MustInit(service, err)

	return conn
}
