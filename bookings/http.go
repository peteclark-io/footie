package bookings

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
	return "bookings"
}

func (h *httpHandler) Create(w http.ResponseWriter, req *http.Request) {
	b, status, err := unmarshalBooking(req.Body)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	status, err = h.repository.Create(b)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(b)
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

	resources.Response(w, http.StatusAccepted, "Booking deleted")
}

func (h *httpHandler) Read(w http.ResponseWriter, req *http.Request) {
	id, status, err := resources.CheckID(tableKey, req)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	b, status, err := h.repository.Read(id)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(b)
}

func (h *httpHandler) Write(w http.ResponseWriter, req *http.Request) {
	id, status, err := resources.CheckID(tableKey, req)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	b, status, err := unmarshalBooking(req.Body)
	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	status, err = h.repository.Write(id, b)

	if err != nil {
		resources.Response(w, status, err.Error())
		return
	}

	resources.Response(w, http.StatusOK, "Saved booking")
}