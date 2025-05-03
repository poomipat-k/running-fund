package operationConfig

import (
	"errors"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

type OperationConfigHandler struct {
	store OperationConfigStore
}

func NewOperationConfigHandler(s OperationConfigStore) *OperationConfigHandler {
	return &OperationConfigHandler{
		store: s,
	}
}

func (h *OperationConfigHandler) GetOperationConfig(w http.ResponseWriter, r *http.Request) {
	conf, err := h.store.GetLatestConfig()
	if err != nil {
		utils.ErrorJSON(w, errors.New("can not load operation config data"), "operationConfig", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, conf)
}
