package v1

import (
	"encoding/json"
	"github.com/anpotashev/mpdgo/internal/api"
	log "github.com/anpotashev/mpdgo/internal/logger"
	"github.com/gorilla/mux"
	"net/http"
)

func (v1 *v1Router) newPlayerRouter(router *mux.Router) {
	router.HandleFunc("/control/{action}", v1.playerControlHandler).Methods(http.MethodPost)
	router.HandleFunc("/playPos", v1.playPosHandler).Methods(http.MethodPost)
	router.HandleFunc("/playId", v1.playIdHandler).Methods(http.MethodPost)
	router.HandleFunc("/seek", v1.seekHandler).Methods(http.MethodPost)
}

func (v1 *v1Router) playIdHandler(w http.ResponseWriter, r *http.Request) {
	var rq PlayIdRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).PlayId(rq.Id), w, r)
}

func (v1 *v1Router) playPosHandler(w http.ResponseWriter, r *http.Request) {
	var rq PlayPosRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).PlayPos(rq.Pos), w, r)
}

func (v1 *v1Router) seekHandler(w http.ResponseWriter, r *http.Request) {
	var rq SeekPosRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).Seek(rq.Pos, rq.SeekTime), w, r)
}

func (v1 *v1Router) playerControlHandler(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]
	log.DebugContext(r.Context(), "player control", "action", action)
	var f func() error
	switch action {
	case "play":
		f = v1.MpdApi.WithRequestContext(r.Context()).Play
	case "pause":
		f = v1.MpdApi.WithRequestContext(r.Context()).Pause
	case "stop":
		f = v1.MpdApi.WithRequestContext(r.Context()).Stop
	case "prev":
		f = v1.MpdApi.WithRequestContext(r.Context()).Previous
	case "next":
		f = v1.MpdApi.WithRequestContext(r.Context()).Next
	default:
		api.ErrorResponse(w, r, http.StatusBadRequest, "Invalid action")
		return
	}
	checkErrorAndWriteResponse(f(), w, r)
}
