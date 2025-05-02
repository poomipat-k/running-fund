package mw

import (
	"errors"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

func AllowCreateNewProject(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowCreateProject := false
		if !allowCreateProject {
			utils.ErrorJSON(w, errors.New("ปิดรับข้อเสนอโครงการ"), "newProjectNotAllow", http.StatusBadRequest)
			return
		}
		next(w, r)
	})
}
