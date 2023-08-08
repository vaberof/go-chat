package middleware

import (
	"context"
	service "github.com/vaberof/go-chat/internal/domain/chat/auth"
	"github.com/vaberof/go-chat/pkg/auth"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"net/http"
)

func AuthMiddleware(next http.Handler, authService service.AuthService, logs *logs.Logs) http.Handler {
	loggerName := "websocket-auth-middleware"
	logger := logs.WithName(loggerName)

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		token := request.URL.Query().Get("token")

		if token == "" {
			logger.Errorw("Client not authenticated: empty token")
			responseWriter.WriteHeader(http.StatusUnauthorized)
			responseWriter.Write([]byte("Need to provide jwt token"))
			return
		}

		userId, err := authService.VerifyToken(token)
		if err != nil {
			logger.Errorf("Invalid token: %v", err)
			// if is.ErrorsAs(ErrTokenExpired) ...
			responseWriter.WriteHeader(http.StatusUnauthorized)
			responseWriter.Write([]byte("Invalid token"))
			return
		}

		logger.Infow("Client is authenticated")

		ctx := auth.UserIdToContext(context.Background(), userId)

		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}
