package api

import (
	"../models"
	"../daos"
	"github.com/gorilla/mux"
	"encoding/json"
	"net/http"
	"github.com/op/go-logging"
	"strconv"
)

const MODULE = "api"

var log = logging.MustGetLogger(MODULE)

type (
	// LeiService specifies the interface for the lei service needed by LeiResource.
	LeiService interface {
		GetAll() ([]models.Lei, error)
		Get(id string) (models.Lei, error)
	}

	// LeiResource defines the handlers for the CRUD APIs.
	LeiResource struct {
		service LeiService
	}
)

// ServeLeiResource sets up the routing of leis endpoints and the corresponding handlers.
func ServeLeiResource(router *mux.Router, service LeiService) {
	resource := &LeiResource{service}
	router.HandleFunc("/leis", resource.getAll)
	router.HandleFunc("/leis/{id:[0-9]+}", resource.get)
}

func (r LeiResource) getAll(writer http.ResponseWriter, request *http.Request) {
	log.Info(request.Proto, request.Host, request.Method, request.RequestURI)

	leis, err := r.service.GetAll()

	switch err {
	case daos.ErrorNoItemFound:
		ItemNotFound(writer, "unable to find data")
	case daos.ErrorDataBaseConnection:
		InternalServerErrorResponse(writer, "unable to connect to database")
	}

	if err != nil {
		message := "unable to establish a connection with the database"
		log.Error(message)
		InternalServerErrorResponse(writer, message)
		return
	}

	payload, err := json.Marshal(leis)
	if err != nil {
		message := "error encoding json"
		log.Error(message)
		InternalServerErrorResponse(writer, message)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}

func (r LeiResource) get(writer http.ResponseWriter, request *http.Request) {
	log.Info(request.Proto, request.Host, request.Method, request.RequestURI)

	vars := mux.Vars(request)
	id := vars["id"]
	lei, err := r.service.Get(id)

	switch err {
	case daos.ErrorNoItemFound:
		ItemNotFound(writer, "unable to find data")
		return
	case daos.ErrorDataBaseConnection:
		InternalServerErrorResponse(writer, "unable to connect to database")
		return
	}

	payload, err := json.Marshal(lei)
	if err != nil {
		message := "error encoding json"
		log.Error(message)
		InternalServerErrorResponse(writer, message)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)

}

func InternalServerErrorResponse(w http.ResponseWriter, message string) {
	status := strconv.Itoa(http.StatusInternalServerError)
	body := map[string]string{"status": status, "message": message}
	payload, _ := json.Marshal(body)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func ItemNotFound(w http.ResponseWriter, message string) {
	status := strconv.Itoa(http.StatusNotFound)
	body := map[string]string{"status": status, "message": message}
	payload, _ := json.Marshal(body)

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
