package matches

import (
	"encoding/json"
	"net/http"

	"github.com/peteclark-io/footie/resources"
)

type httpHandler struct {
	repository *Repository
}

func NewHTTPHandler() resources.Handler {
	return &httpHandler{repository: &Repository{}}
}

func (h *httpHandler) Name() string {
	return "matches"
}

func (h *httpHandler) Create(w http.ResponseWriter, req *http.Request) {
	m, status, err := unmarshalMatch(req.Body)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	status, err = h.repository.Create(m)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.AddCommonHeaders(w)
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(m)
}

func (h *httpHandler) Delete(w http.ResponseWriter, req *http.Request) {
	id, status, err := resources.CheckID(tableKey, req)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	status, err = h.repository.Delete(id)

	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.Response(w, http.StatusAccepted, "Match deleted")
}

func (h *httpHandler) Read(w http.ResponseWriter, req *http.Request) {
	id, status, err := resources.CheckID(tableKey, req)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	m, status, err := h.repository.Read(id)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.AddCommonHeaders(w)
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(m)
}

func (h *httpHandler) Write(w http.ResponseWriter, req *http.Request) {
	id, status, err := resources.CheckID(tableKey, req)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	m, status, err := unmarshalMatch(req.Body)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	status, err = h.repository.Write(id, m)

	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.Response(w, http.StatusOK, "Saved match")
}
