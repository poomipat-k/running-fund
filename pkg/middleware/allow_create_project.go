package mw

import (
	"errors"
	"net/http"

	operationConfig "github.com/poomipat-k/running-fund/pkg/operation-config"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

func AllowCreateNewProject(next http.HandlerFunc, confStore operationConfig.OperationConfigStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf, err := confStore.GetLatestConfig()
		if err != nil {
			utils.ErrorJSON(w, errors.New("can not load operation config data"), "operationConfig", http.StatusInternalServerError)
			return
		}

		if !*conf.AllowNewProject {
			utils.ErrorJSON(w, errors.New("ปิดรับข้อเสนอโครงการ"), "newProjectNotAllow", http.StatusBadRequest)
			return
		}
		next(w, r)
	})
}
