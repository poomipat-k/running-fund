package address

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

type AddressStore interface {
	GetProvinces() ([]Province, error)
}

type AddressHandler struct {
	store AddressStore
}

func NewAddressHandler(s AddressStore) *AddressHandler {
	return &AddressHandler{
		store: s,
	}
}

func (h *AddressHandler) GetProvinces(w http.ResponseWriter, r *http.Request) {
	log.Println("GetProvinces Handler")
	provinces, err := h.store.GetProvinces()
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "province")
		return
	}
	utils.WriteJSON(w, http.StatusOK, provinces)
}
