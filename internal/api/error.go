package api

import (
	"github.com/anpotashev/mpdgo/internal/api/middleware"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, code int, msg string) {
	requestId, success := r.Context().Value(middleware.RequestIdContextAttributeName).(string)
	var err error
	if !success {
		err = WriteJson(w, code, Envelope{"success": false, "error": msg})
	} else {
		err = WriteJson(w, code, Envelope{"success": false, "error": msg, "requestId": requestId})
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
