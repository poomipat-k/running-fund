package utils

import (
	"net/http"
	"strconv"
)

func GetUserIdFromRequestHeader(r *http.Request) (int, error) {
	strUserId := r.Header.Get("userId")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func GetUserRoleFromRequestHeader(r *http.Request) string {
	return r.Header.Get("userRole")
}
