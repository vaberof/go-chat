package http

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
)

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *loginRequestBody) Bind(r *http.Request) error {
	return nil
}

type loginResponsePayload struct {
	Token string `json:"token"`
}

func (h *Handler) Login() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginReqBody := &loginRequestBody{}
		if err := render.Bind(r, loginReqBody); err != nil {
			views.RenderJSON(w, r, http.StatusBadRequest, apiv1.Error(InvalidRequestBodyMessage))
			return
		}

		token, err := h.authService.Login(loginReqBody.Username, loginReqBody.Password)
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(loginResponsePayload{Token: string(*token)})

		views.RenderJSON(w, r, http.StatusOK, apiv1.Success(payload))
	})
}
