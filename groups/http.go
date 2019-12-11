package groups

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
	return groupsResource
}

func (h *httpHandler) Create(w http.ResponseWriter, req *http.Request) {
	gr, status, err := unmarshalGroup(req.Body)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	status, err = h.repository.Create(gr)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.AddCommonHeaders(w)
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(gr)
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

	resources.Response(w, http.StatusAccepted, "Group deleted")
}

func (h *httpHandler) Read(w http.ResponseWriter, req *http.Request) {
	id, status, err := resources.CheckID(tableKey, req)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	gr, status, err := h.repository.Read(id)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.AddCommonHeaders(w)
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.Encode(gr)
}

func (h *httpHandler) Write(w http.ResponseWriter, req *http.Request) {
	id, status, err := resources.CheckID(tableKey, req)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	gr, status, err := unmarshalGroup(req.Body)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	status, err = h.repository.Write(id, gr)

	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.Response(w, http.StatusOK, "Saved group")
}
