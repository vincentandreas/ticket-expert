package utilities

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"ticket-expert/models"
)

func HashParams(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}

func WriteErrorResp(w http.ResponseWriter, httpStatus int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	temp := models.ApiResponse{"", "01", errorMsg}
	json.NewEncoder(w).Encode(temp)
}

func WriteSuccessResp(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	temp := models.ApiResponse{"", "00", "Success"}
	json.NewEncoder(w).Encode(temp)
}
