package chiworker

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/voi-oss/svc"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"strconv"
)

const (
	DefaultHTTPPort = 8080
)

type Controller interface {
	Init(*zap.Logger) error
	SetupRouter(router chi.Router) error
	Terminate() error
}

var _ svc.Worker = (*Worker)(nil)

type Worker struct {
	name   string
	prefix string
	port   int
	logger *zap.Logger
	router chi.Router
	ctrl   Controller
	server *http.Server
}

func New(c Controller) *Worker {
	w := &Worker{
		name:   "chi-http",
		prefix: "",
		port:   DefaultHTTPPort,
		ctrl:   c,
	}
	return w
}

func (w Worker) Init(logger *zap.Logger) error {
	w.logger = logger
	w.router = chi.NewRouter()
	httpLogger, err := zap.NewStdLogAt(w.logger, zapcore.ErrorLevel)
	if err != nil {
		return err
	}
	w.server = &http.Server{
		Addr: net.JoinHostPort("", strconv.Itoa(w.port)),
		Handler: &ochttp.Handler{
			Handler: w.router,
		},
		ErrorLog: httpLogger,
	}

	if err := w.ctrl.Init(w.logger); err != nil {
		return fmt.Errorf("init error for controller, error:%w", err)
	}

	r := chi.NewRouter()
	if err := w.ctrl.SetupRouter(r); err != nil {
		return fmt.Errorf("SetupRouter error for controller, error:%w", err)
	}
	w.router.Mount(w.prefix, r)

	return nil
}

func (w *Worker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w.router.ServeHTTP(rw, req)
}

func (w *Worker) Name() string {
	return w.name
}

func (w Worker) Run() error {
	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (w Worker) Terminate() error {
	return w.ctrl.Terminate()
}
