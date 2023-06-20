package logging

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const reqLogMsg = "transport_log"

func (l *Logger) NewHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		next.ServeHTTP(w, r)

		l.Info(
			reqLogMsg,
			String("method", r.Method),
			String("path", r.URL.Path),
			String("duration", time.Since(begin).String()),
		)
	})
}

func (l *Logger) NewEchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		begin := time.Now()

		err := next(c)

		if err != nil {
			l.Error(
				reqLogMsg,
				String("method", c.Request().Method),
				String("path", c.Request().URL.Path),
				String("duration", time.Since(begin).String()),
				String("error", err.Error()),
			)
			return err
		}

		l.Info(
			reqLogMsg,
			String("method", c.Request().Method),
			String("path", c.Request().URL.Path),
			String("duration", time.Since(begin).String()),
		)

		return err
	}
}
