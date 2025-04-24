package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int `json:"code"`
	Token  *string `json:"token,omitempty"`
	Total int `json:"total,omitempty"`
	Message string `json:"message"`
	Data    any `json:"data"`
}

func SendResponse(w http.ResponseWriter, token string, codes int, msg string, data any) {
	w.WriteHeader(codes)
	Resp := Response{
		Code: codes,
		Message: msg,
		Data: data,
	}
	if token != "" {
		Resp.Token = &token
	}
	if err := json.NewEncoder(w).Encode(Resp); err != nil {
		http.Error(w, "Failed to encode response message", http.StatusInternalServerError)
	}
}

func SendList(w http.ResponseWriter, codes, total int, msg string, data any) {
	w.WriteHeader(codes)
	Resp := Response{
		Code: codes,
		Total: total,
		Message: msg,
		Data: data,
	}
	if err := json.NewEncoder(w).Encode(Resp); err != nil {
		http.Error(w, "Failed to encode response list", http.StatusInternalServerError)
	}
}