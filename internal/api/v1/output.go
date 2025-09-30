package v1

import (
	"encoding/json"
	"github.com/anpotashev/mpdgo/internal/api"
	"github.com/anpotashev/mpdgo/internal/api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (v1 *v1Router) newOutputRouter(router *mux.Router) {
	router.HandleFunc("/list", v1.getOutputsHandler).Methods(http.MethodGet)
	router.HandleFunc("/enable", v1.enableOutputHandler).Methods(http.MethodPost)
	router.HandleFunc("/disable", v1.disableOutputHandler).Methods(http.MethodPost)
}

func (v1 *v1Router) enableOutputHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(w, r)
	if err != nil {
		return
	}
	v1.setOutputState(w, r, id, true)
}

func (v1 *v1Router) disableOutputHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(w, r)
	if err != nil {
		return
	}
	v1.setOutputState(w, r, id, false)
}

func getIdFromRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	var rq ChangeOutputStateRequest
	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return 0, err
	}
	return rq.OutputId, nil
}

func (v1 *v1Router) setOutputState(w http.ResponseWriter, r *http.Request, outputId int, enable bool) {
	var f func(int) error
	if enable {
		f = v1.MpdApi.WithRequestContext(r.Context()).EnableOutput
	} else {
		f = v1.MpdApi.WithRequestContext(r.Context()).DisableOutput
	}
	checkErrorAndWriteResponse(f(outputId), w, r)
}

func (v1 *v1Router) getOutputsHandler(w http.ResponseWriter, r *http.Request) {
	outputs, err := v1.MpdApi.WithRequestContext(r.Context()).ListOutputs()
	var payload interface{}
	if err == nil {
		payload = dto.MapSlice(outputs, dto.MapOutput)
	}
	checkErrorAndWriteResponseWithPayload(payload, err, w, r)
}
