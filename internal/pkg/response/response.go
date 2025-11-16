package response

import (
	"encoding/json"
	"net/http"

	"service-pr-reviewer-assignment/internal/generated/api/dto"
)

func OK(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func Error(w http.ResponseWriter, status int, errorCode dto.ErrorResponseErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(
		dto.ErrorResponse{
			Error: struct {
				Code    dto.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    errorCode,
				Message: message,
			},
		},
	)
}
