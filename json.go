package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func DecodeReq[T any](r *http.Request) (T, error) {
	var payload T

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func DecodeJSON[T any](r io.Reader) (T, error) {
	var payload T

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	err := dec.Decode(&payload)
	if err != nil {
		return payload, fmt.Errorf("decode error: %v", err)
	}

	return payload, nil
}

func EncodeJSON[T any](w http.ResponseWriter, payload *T) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(payload)
}

func JsonError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

type response[T any] struct {
	Error   string `json:"error"`
	Payload *T     `json:"payload"`
}
