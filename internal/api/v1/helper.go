package v1

import (
	"github.com/anpotashev/mpdgo/internal/api"
	"net/http"
)

func checkErrorAndWriteResponse(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	api.WriteJson(w, http.StatusOK, api.Envelope{
		"success": true,
	})
}

func checkErrorAndWriteResponseWithPayload(payload interface{}, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	api.WriteJson(w, http.StatusOK, api.Envelope{
		"success": true,
		"payload": payload,
	})
}
