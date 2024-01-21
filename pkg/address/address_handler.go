package address

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

type AddressStore interface {
	GetProvinces() ([]Province, error)
	GetDistrictsByProvince(provinceId int) ([]District, error)
	GetSubdistrictsByDistrict(districtId int) ([]Subdistrict, error)
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
	provinces, err := h.store.GetProvinces()
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "province")
		return
	}
	utils.WriteJSON(w, http.StatusOK, provinces)
}

func (h *AddressHandler) GetDistrictsByProvince(w http.ResponseWriter, r *http.Request) {
	p := chi.URLParam(r, "provinceId")
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

func (h *AddressHandler) GetSubdistrictsByProvince(w http.ResponseWriter, r *http.Request) {
	p := chi.URLParam(r, "districtId")
	if p == "" {
		utils.ErrorJSON(w, errors.New("districtId is required"), "districtId")
		return
	}
	districtId, err := strconv.Atoi(p)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "districtId")
		return
	}
	subdistricts, err := h.store.GetSubdistrictsByDistrict(districtId)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}
	utils.WriteJSON(w, http.StatusOK, subdistricts)
}
