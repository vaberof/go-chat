package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/internal/domain/chat/auth"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
)

type registerRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *registerRequestBody) Bind(req *http.Request) error {
	return nil
}

type registerResponsePayload struct {
	UserId   domain.UserId `json:"user_id"`
	Username string        `json:"username"`
}

func RegisterRoute(authService auth.AuthService) func(router chi.Router) {
	return func(router chi.Router) {
		router.Post("/register", RegisterHandler(authService))
	}
}

func RegisterHandler(authService auth.AuthService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registerReqBody := &registerRequestBody{}
		if err := render.Bind(r, registerReqBody); err != nil {
			views.RenderJSON(w, r, http.StatusBadRequest, apiv1.Error(InvalidRequestBodyMessage))
			return
		}

		user, err := authService.Register(registerReqBody.Username, registerReqBody.Password)
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(registerResponsePayload{UserId: user.Id, Username: user.Username})

		views.RenderJSON(w, r, http.StatusOK, apiv1.Success(payload))
	})
}
