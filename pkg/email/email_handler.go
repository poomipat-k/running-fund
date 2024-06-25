package appEmail

import "net/http"

type EmailHandler struct{}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{}
}

func (h *EmailHandler) HandlingBounces(w http.ResponseWriter, r *http.Request) {

}
