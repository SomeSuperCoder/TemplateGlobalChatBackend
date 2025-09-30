package utils

import (
	"fmt"
	"net/http"
)

func CheckError(w http.ResponseWriter, err error, message string, code int) bool {
	if err != nil {
		http.Error(w, fmt.Sprintf("message: %v", err.Error()), code)
		return true
	}
	return false
}
