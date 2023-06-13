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
	SetAllHeaders(w)
	w.WriteHeader(httpStatus)
	temp := models.ApiResponse{"", "01", errorMsg}
	json.NewEncoder(w).Encode(temp)
}

func WriteSuccessResp(w http.ResponseWriter) {
	SetAllHeaders(w)
	w.WriteHeader(http.StatusOK)
	temp := models.ApiResponse{"", "00", "Success"}
	json.NewEncoder(w).Encode(temp)
}

func WriteUnauthResp(w http.ResponseWriter) {
	SetAllHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)
	temp := models.ApiResponse{"", "01", "Unauthorized"}
	json.NewEncoder(w).Encode(temp)
}

func WriteSuccessWithDataResp(w http.ResponseWriter, data interface{}) {
	SetAllHeaders(w)

	w.WriteHeader(http.StatusOK)
	temp := models.ApiGetResponse{data, "00", "Success"}
	json.NewEncoder(w).Encode(temp)
}

func SetAllHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}
