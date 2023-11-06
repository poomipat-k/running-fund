package server

import (
	"encoding/json"
	"net/http"
)

func ResposeJson(w http.ResponseWriter, body any, httpStatus int) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if httpStatus == 0 {
		httpStatus = http.StatusOK
	}
	w.WriteHeader(httpStatus)
	w.Write(jsonBytes)
}
