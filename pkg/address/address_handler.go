package address

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

type AddressStore interface {
	GetProvinces() ([]Province, error)
	GetDistrictsByProvince(provinceId int) ([]District, error)
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

func (h *AddressHandler) GetDistrictsByProvince(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDistrictsByProvinceRequest Handler")
	p := chi.URLParam(r, "provinceId")
	log.Println("===p", p)
	if p == "" {
		utils.ErrorJSON(w, errors.New("provinceId is required"), "provinceId")
		return
	}
	provinceId, err := strconv.Atoi(p)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "provinceId")
		return
	}
	districts, err := h.store.GetDistrictsByProvince(provinceId)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}
	utils.WriteJSON(w, http.StatusOK, districts)

}
