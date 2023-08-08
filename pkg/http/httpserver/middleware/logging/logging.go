package logging

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"net/http"
)

type Middleware struct {
	Handler func(handler http.Handler) http.Handler
	Logger  *zap.SugaredLogger
}

func New(logs *logs.Logs) *Middleware {
	return impl(logs, "")
}

func impl(logs *logs.Logs, serverName string) *Middleware {
	loggerName := "http-server"
	if serverName != "" {
		loggerName = fmt.Sprintf("%s.%s", loggerName, serverName)
	}
	logger := logs.WithName(loggerName)

	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			path := request.URL.Path
			if path == "" {
				path = "/"
			}
			method := request.Method
			logger.Infow("Request started", "http.path", path, "http.method", method)

			ww := middleware.NewWrapResponseWriter(responseWriter, request.ProtoMajor)

			defer func() {
				status := ww.Status()
				if status == 0 {
					s, ok := request.Context().Value(render.StatusCtxKey).(int)
					if ok && s != status {
						status = s
					}
				}

				if status >= 500 {
					logger.Infow("Request finished", "http.path", path, "http.method", method, "http.result", "error", "http.status", status)
					return
				}

				logger.Infow("Request finished", "http.path", path, "http.method", method, "http.result", "success", "http.status", status)
			}()

			next.ServeHTTP(ww, request)
		})
	}

	return &Middleware{
		Handler: handler,
		Logger:  logger,
	}
}
