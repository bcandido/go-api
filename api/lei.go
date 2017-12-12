package api

import (
	"../models"
	"../service"
	"github.com/gorilla/mux"
	"encoding/json"
	"net/http"
	"github.com/op/go-logging"
	"strconv"
	validate "github.com/gima/govalid/v1"
	"fmt"
)

const MODULE = "api"
const itemNotFound = "item nor found"
const internalServerError = "internal server error"
const jsonDecodingError = "error decoding json"
const jsonEncodingError = "error encoding json"
const bodyBadFormatted = "body bad formatted"

var log = logging.MustGetLogger(MODULE)

type (
	// LeiService specifies the interface for the lei service needed by LeiResource.
	LeiService interface {
		GetAll() ([]models.Lei, error)
		Get(id string) (models.Lei, error)
		Add(name string) (error)
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
	if err = service.Validate(err); err != nil {
		log.Error(err.Error())
		switch err {
		case service.ErrorNoItemFound:
			ItemNotFoundResponse(writer, itemNotFound)
		default:
			InternalServerErrorResponse(writer, internalServerError)
		}
		return
	}

	payload, err := json.Marshal(leis)
	if err != nil {
		log.Error(err.Error())
		InternalServerErrorResponse(writer, jsonEncodingError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}

func (r LeiResource) get(writer http.ResponseWriter, request *http.Request) {
	log.Info(request.Proto, request.Host, request.Method, request.RequestURI)

	// get lei id from request
	vars := mux.Vars(request)
	id := vars["id"]

	lei, err := r.service.Get(id)
	if err = service.Validate(err); err != nil {
		log.Error(err.Error())
		switch err {
		case service.ErrorNoItemFound:
			ItemNotFoundResponse(writer, err.Error())
		default:
			InternalServerErrorResponse(writer, err.Error())
		}
		return
	}

	payload, err := json.Marshal(lei)
	if err != nil {
		log.Error(err.Error())
		InternalServerErrorResponse(writer, jsonDecodingError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)

}

func (r LeiResource) add(writer http.ResponseWriter, request *http.Request) {
	log.Info(request.Proto, request.Host, request.Method, request.RequestURI)

	decoder := json.NewDecoder(request.Body)
	var body = make(map[string]string)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		InternalServerErrorResponse(writer, jsonDecodingError)
		return
	}

	if err = ValidateJson(body); err != nil {
		log.Error(err.Error())
		InternalServerErrorResponse(writer, bodyBadFormatted)
		return
	}

	err = r.service.Add(body["nome"])
	if err = service.Validate(err); err != nil {
		log.Error(err.Error())
		switch err {
		case service.ErrorNoItemFound:
			ItemNotFoundResponse(writer, itemNotFound)
		case service.ErrorAlreadyInserted:
			InternalServerErrorResponse(writer, "lei já inserida")
		default:
			InternalServerErrorResponse(writer, internalServerError)
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")

	payload, _ := json.Marshal(map[string]string{"status": string(200), "message": "OK"})
	writer.Write(payload)
}

func ValidateJson(data map[string]string) error {
	schema := validate.Object(validate.ObjKV("nome", validate.String()))
	if path, err := schema.Validate(data); err == nil {
		log.Info("validation passed.")
		return nil
	} else {
		log.Error(fmt.Sprintf("validation failed at %s. Error (%s)", path, err))
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

func ItemNotFoundResponse(w http.ResponseWriter, message string) {
	status := strconv.Itoa(http.StatusNotFound)
	body := map[string]string{"status": status, "message": message}
	payload, _ := json.Marshal(body)

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
