package v1

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (v1 *v1Router) newConnectionRouter(router *mux.Router) {
	router.HandleFunc("/connect", v1.connectHandler).Methods(http.MethodPost)
	router.HandleFunc("/disconnect", v1.disconnectHandler).Methods(http.MethodPost)
	router.HandleFunc("/state", v1.getConnectionState).Methods(http.MethodGet)
}

func (v1 *v1Router) connectHandler(w http.ResponseWriter, r *http.Request) {
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).Connect(), w, r)
}

func (v1 *v1Router) disconnectHandler(w http.ResponseWriter, r *http.Request) {
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).Disconnect(), w, r)
}

func (v1 *v1Router) getConnectionState(w http.ResponseWriter, r *http.Request) {
	result := v1.MpdApi.WithRequestContext(r.Context()).IsConnected()
	checkErrorAndWriteResponseWithPayload(result, nil, w, r)
}
