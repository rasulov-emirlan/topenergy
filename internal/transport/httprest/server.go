package httprest

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"

	"github.com/rasulov-emirlan/topenergy-interview/config"
	"github.com/rasulov-emirlan/topenergy-interview/internal/domains"
	"github.com/rasulov-emirlan/topenergy-interview/pkg/health"
	"github.com/rasulov-emirlan/topenergy-interview/pkg/logging"
)

const ServiceName = "tasks-service"

type server struct {
	srv *http.Server
}

func NewServer(cfg config.Config) server {
	return server{
		srv: &http.Server{
			Addr:         cfg.Server.Port,
			ReadTimeout:  cfg.Server.TimeoutRead,
			WriteTimeout: cfg.Server.TimeoutWrite,
		},
	}
}

type validatorWrapper struct {
	validator *validator.Validate
}

func (v *validatorWrapper) Validate(i any) error {
	return v.validator.Struct(i)
}

func (s server) Start(log *logging.Logger, doms domains.DomainCombiner, checks []health.Checker) error {
	router := echo.New()
	router.HideBanner = true
	router.HidePort = true
	router.Validator = &validatorWrapper{validator: validator.New()}
	router.Use(log.NewEchoMiddleware)
	router.Use(middleware.Gzip())
	router.Use(middleware.Recover())
	router.Use(middleware.CORS())

	router.Use(otelecho.Middleware(ServiceName))
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		ctx := c.Request().Context()
		trace.SpanFromContext(ctx).RecordError(err)

		router.DefaultHTTPErrorHandler(err, c)
	}

	router.Any("/health*", echo.WrapHandler(health.NewHTTPHandler(ServiceName, checks)))
	router.GET("/routes", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, router.Routes())
	})

	tasksHandler := NewTasksHandler(doms.TasksService())
	tasksGroup := router.Group("/tasks")
	{
		tasksGroup.POST("", tasksHandler.Create)
		tasksGroup.GET("", tasksHandler.ReadAll)
		tasksGroup.GET("/:id", tasksHandler.Read)
		tasksGroup.PUT("/:id", tasksHandler.Update)
		tasksGroup.DELETE("/:id", tasksHandler.Delete)
	}

	s.srv.Handler = router
	return s.srv.ListenAndServe()
}

func (s server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
