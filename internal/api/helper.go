package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Envelope map[string]any

func WriteJson(w http.ResponseWriter, statusCode int, data Envelope) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func ReadJson(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	return decoder.Decode(dst)
}

func GetIntValue(r *http.Request, key string) (int, error) {
	return strconv.Atoi(r.URL.Query().Get(key))
}

func WriteError(w http.ResponseWriter, statusCode int, msg string) {
	WriteJson(w, statusCode, Envelope{"error": msg})
}
