package v1

import (
	"github.com/anpotashev/mpdgo/internal/api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (v1 *v1Router) newTreeRouter(router *mux.Router) {
	router.HandleFunc("/get", v1.getTreeHandler).Methods(http.MethodGet)
}

func (v1 *v1Router) getTreeHandler(w http.ResponseWriter, r *http.Request) {
	tree, err := v1.MpdApi.WithRequestContext(r.Context()).Tree()
	var payload interface{}
	if err == nil {
		payload = dto.MapMpdTree(*tree)
	}
	checkErrorAndWriteResponseWithPayload(payload, err, w, r)
}
