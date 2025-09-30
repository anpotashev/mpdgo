package v1

import (
	"github.com/anpotashev/mpdgo/internal/api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (v1 *v1Router) newStatusRouter(router *mux.Router) {
	router.HandleFunc("/get", v1.getStatusHandler).Methods(http.MethodGet)
}

func (v1 *v1Router) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := v1.MpdApi.WithRequestContext(r.Context()).Status()
	var payload interface{}
	if err == nil {
		payload = dto.MapStatus(status)
	}
	checkErrorAndWriteResponseWithPayload(payload, err, w, r)
}
