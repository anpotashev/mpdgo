package v1

import (
	"encoding/json"
	"fmt"
	"github.com/anpotashev/mpdgo/internal/api"
	"github.com/anpotashev/mpdgo/internal/api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (v1 *v1Router) newPlaylistRouter(router *mux.Router) {
	router.HandleFunc("/list", v1.getPlaylistsHandler).Methods(http.MethodGet)
	router.HandleFunc("/info", v1.getPlaylistInfoHandler).Methods(http.MethodPost)
	router.HandleFunc("/delete", v1.deletePlaylistHandler).Methods(http.MethodDelete)
	router.HandleFunc("/rename", v1.renamePlaylistHandler).Methods(http.MethodPost)
	router.HandleFunc("/save", v1.saveCurrentPlaylistAsStored).Methods(http.MethodPost)
}

func (v1 *v1Router) getPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	playlists, err := v1.MpdApi.WithRequestContext(r.Context()).GetPlaylists()
	var payload []dto.Playlist
	if err == nil {
		fmt.Println("here")
		payload = dto.MapSlice(playlists, dto.MapPlaylist)
	}
	checkErrorAndWriteResponseWithPayload(payload, err, w, r)
}

func (v1 *v1Router) getPlaylistInfoHandler(w http.ResponseWriter, r *http.Request) {
	var request GetStoredPlaylistsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	playlist, err := v1.MpdApi.WithRequestContext(r.Context()).PlaylistInfo(request.Name)
	var payload dto.Playlist
	if err == nil {
		payload = dto.MapPlaylist(*playlist)
	}
	checkErrorAndWriteResponseWithPayload(payload, err, w, r)
}

func (v1 *v1Router) deletePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	var request DeleteStoredPlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).DeleteStoredPlaylist(request.Name), w, r)
}

func (v1 *v1Router) renamePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	var request RenameStoredPlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).RenameStoredPlaylist(request.OldName, request.NewName), w, r)
}

func (v1 *v1Router) saveCurrentPlaylistAsStored(w http.ResponseWriter, r *http.Request) {
	var request SaveCurrentPlaylistAsStoredRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	checkErrorAndWriteResponse(v1.MpdApi.WithRequestContext(r.Context()).SaveCurrentPlaylistAsStored(request.Name), w, r)
}
