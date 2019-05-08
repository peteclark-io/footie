package resources

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/husobee/vestigo"
)

type Msg struct {
	Msg string `json:"message"`
}

type Handler interface {
	Name() string
	Read(w http.ResponseWriter, r *http.Request)
	Write(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func Response(w http.ResponseWriter, status int, msg string) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(&Msg{Msg: msg})
}

func CheckID(tableKey string, req *http.Request) (string, int, error) {
	id := vestigo.Param(req, tableKey)
	if id == "" {
		return id, http.StatusBadRequest, errors.New("Please provide an 'id'")
	}
	return id, 0, nil
}
