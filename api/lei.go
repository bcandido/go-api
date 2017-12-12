package api

import (
	"../models"
	"../daos"
	"github.com/gorilla/mux"
	"encoding/json"
	"net/http"
	"github.com/op/go-logging"
	"strconv"
	"io/ioutil"
	validate "github.com/gima/govalid/v1"
	"fmt"
)

const MODULE = "api"

var log = logging.MustGetLogger(MODULE)

type (
	// LeiService specifies the interface for the lei service needed by LeiResource.
	LeiService interface {
		GetAll() ([]models.Lei, error)
		Get(id string) (models.Lei, error)
		Add(name string) (bool, error)
	}

	// LeiResource defines the handlers for the CRUD APIs.
	LeiResource struct {
		service LeiService
	}
)

// ServeLeiResource sets up the routing of leis endpoints and the corresponding handlers.
func ServeLeiResource(router *mux.Router, service LeiService) {
	resource := &LeiResource{service}
	router.HandleFunc("/leis", resource.getAll).Methods("GET")
	router.HandleFunc("/leis", resource.add).Methods("POST")
	router.HandleFunc("/leis/{id:[0-9]+}", resource.get).Methods("GET")
}

func (r LeiResource) getAll(writer http.ResponseWriter, request *http.Request) {
	log.Info(request.Proto, request.Host, request.Method, request.RequestURI)

	leis, err := r.service.GetAll()

	switch err {
	case daos.ErrorNoItemFound:
		ItemNotFound(writer, err.Error())
		return
	case daos.ErrorDataBaseConnection:
		InternalServerErrorResponse(writer, err.Error())
		return
	case daos.ErrorTransactionFailure:
		InternalServerErrorResponse(writer, err.Error())
		return
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
		ItemNotFound(writer, err.Error())
		return
	case daos.ErrorDataBaseConnection:
		InternalServerErrorResponse(writer, err.Error())
		return
	case daos.ErrorTransactionFailure:
		InternalServerErrorResponse(writer, err.Error())
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

func (r LeiResource) add(writer http.ResponseWriter, request *http.Request) {
	log.Info(request.Proto, request.Host, request.Method, request.RequestURI)

	body, err := ioutil.ReadAll(request.Body)
	// err = validateJson(body)
	if err != nil {
		message := "body has an invalid json"
		log.Error(message)
		InternalServerErrorResponse(writer, message)
		return
	}

	var lei = models.Lei{}
	if err := json.Unmarshal(body, &lei); err != nil {
		message := "error decoding json"
		log.Error(message)
		InternalServerErrorResponse(writer, message)
		return
	}

	r.service.Add(lei.Nome)
	payload, _ := json.Marshal(map[string]string{"status": string(200), "message": "OK"})

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}

func validateJson(data []byte) error {
	schema := validate.Object(validate.ObjKV("nome", validate.String()))
	if path, err := schema.Validate(data); err == nil {
		log.Info("Validation passed.")
		return nil
	} else {
		log.Error(fmt.Sprintf("Validation failed at %s. Error (%s)", path, err))
		return err
	}
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
