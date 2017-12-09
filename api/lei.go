package api

import (
	"../models"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/op/go-logging"
)

const MODULE = "api"

var log = logging.MustGetLogger(MODULE)

type (
	// LeiService specifies the interface for the lei service needed by LeiResource.
	LeiService interface {
		GetAll() ([]models.Lei, error)
	}

	// artistResource defines the handlers for the CRUD APIs.
	LeiResource struct {
		service LeiService
	}
)

// ServeArtist sets up the routing of artist endpoints and the corresponding handlers.
func ServeLeiResource(router *mux.Router, service LeiService) {
	resource := &LeiResource{service}
	router.HandleFunc("/leis", resource.getAll)
}

func (r LeiResource) getAll(writer http.ResponseWriter, request *http.Request) {
	log.Info(request.Proto, request.Host, request.Method, request.RequestURI)

	leis, err := r.service.GetAll()
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
	fmt.Fprintf(writer, string(payload))
}

func InternalServerErrorResponse(w http.ResponseWriter, message string) {
	status := http.StatusInternalServerError
	body := map[int]string{int(status): message}
	payload, _ := json.Marshal(body)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(payload))
}
