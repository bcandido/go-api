package main

import (
	"fmt"
	"github.com/op/go-logging"
	"./app"
	"./db"
	"./dao"
	"net/http"
	"html"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
)

var log = logging.MustGetLogger("main")

func main() {
	log.Info("start application")
	log.Info("Version: " + app.Version)

	log.Info("set up application")
	SetUpApplication()

	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/leis", ServeLeis)

	port := fmt.Sprintf(":%v", app.Config.ServerPort)
	log.Info("server version " + app.Version + " is started at " + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Proto, r.Host, r.Method, r.RequestURI)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func ServeLeis(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Proto, r.Host, r.Method, r.RequestURI)

	// open db connection
	database := db.Postgres{}
	err := database.Open()
	defer database.Close()

	if err != nil {
		message := "unable to establish a connection with the database"
		log.Error(message)
		InternalServerErrorResponse(w, message)
		return
	}

	lei := dao.NewLeiDAO(database.DB)
	leis, err := lei.GetAll()

	response, err := json.Marshal(leis)
	if err != nil {
		message := "error encoding json"
		log.Error(message)
		InternalServerErrorResponse(w, message)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(response))
}

func SetUpApplication() {
	// simply prints an ascii art
	PrintAsciiArt()

	// load application configurations
	LoadConfiguration()
}

func PrintAsciiArt() {
	asciiArt, err := ioutil.ReadFile("ascii-art.txt")
	if err != nil {
		log.Warning("unable to load ascii-art message")
	} else {
		log.Info("\n" + string(asciiArt))
	}
}

func LoadConfiguration() {
	log.Info("load application configurations")
	err := app.LoadConfig("./")
	if err != nil {
		log.Error("invalid application configuration")
		panic(err)
	}
}

func InternalServerErrorResponse(w http.ResponseWriter, message string)  {
	status := http.StatusInternalServerError
	body := map[int]string{int(status): message}

	response, _ := json.Marshal(body)
	w.WriteHeader(status)
	fmt.Fprintf(w, string(response))
}
