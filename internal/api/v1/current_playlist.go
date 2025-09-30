package v1

import (
	"encoding/json"
	"github.com/anpotashev/mpdgo/internal/api"
	"github.com/anpotashev/mpdgo/internal/api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (v1 *v1Router) newCurrentPlaylistRouter(router *mux.Router) {
	router.HandleFunc("/list", v1.listCurrentPlaylistHandler).Methods(http.MethodGet)
	router.HandleFunc("/add", v1.addToCurrentPlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/clear", v1.clearCurrentPlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/addToPos", v1.addToPosToCurrentPlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/deleteByPos", v1.deleteByPosFromCurrentPlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/move", v1.moveCurrentPlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/batchMove", v1.batchMoveCurrentPlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/shuffleAll", v1.shuffleAllCurrentPlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/shuffle", v1.shuffleCurrentPlaylistHandler).Methods(http.MethodPost)

}

func (v1 *v1Router) listCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	playlist, err := v1.MpdApi.WithRequestContext(r.Context()).Playlist()
	var payload interface{}
	if err == nil {
		payload = dto.MapPlaylist(*playlist)
	}
	checkErrorAndWriteResponseWithPayload(payload, err, w, r)
}

func (v1 *v1Router) addToCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	rq := AddToPlaylistRequest{}
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).Add(rq.Path), w, r)
}

func (v1 *v1Router) clearCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).Clear(), w, r)
}

func (v1 *v1Router) addToPosToCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	rq := AddToPosPlaylistRequest{}
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).AddToPos(rq.Pos, rq.Path), w, r)
}

func (v1 *v1Router) deleteByPosFromCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	var rq DeleteByPosFromCurrentPlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).DeleteByPos(rq.Pos), w, r)
}

func (v1 *v1Router) moveCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	var rq MoveInCurrentPlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).Move(rq.FromPos, rq.ToPos), w, r)
}

func (v1 *v1Router) batchMoveCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	var rq BatchMoveInCurrentPlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).BatchMove(rq.FromStartPos, rq.FromEndPos, rq.ToPos), w, r)
}

func (v1 *v1Router) shuffleAllCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).ShuffleAll(), w, r)
}

func (v1 *v1Router) shuffleCurrentPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	var rq ShuffleCurrentPlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).Shuffle(rq.FromPos, rq.ToPos), w, r)
}
