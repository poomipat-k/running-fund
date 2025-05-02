package operationConfig

import (
	"log/slog"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

type OperationConfigStore interface {
	GetLatestConfig() (OperationConfig, error)
	// Update(conf OperationConfig) (bool, error)
}

type OperationConfigHandler struct {
	store OperationConfigStore
}

func NewOperationConfigHandler(s OperationConfigStore) *OperationConfigHandler {
	return &OperationConfigHandler{
		store: s,
	}
}

func (h *OperationConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	conf, err := h.store.GetLatestConfig()
	if err != nil {
		fail(w, err, "operationConfig")
		return
	}
	utils.WriteJSON(w, http.StatusOK, conf)
}

func fail(w http.ResponseWriter, err error, name string, status ...int) {
	slog.Error(err.Error())
	s := http.StatusBadRequest
	if len(status) > 0 {
		s = status[0]
	}
	utils.ErrorJSON(w, err, name, s)
}
