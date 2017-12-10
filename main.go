package main

import (
	"fmt"
	"github.com/op/go-logging"
	"./app"
	"./db"
	"./daos"
	"./api"
	"./service"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
)

var log = logging.MustGetLogger("main")

func main() {
	log.Info("start application version " + app.Version)
	SetUpApplication()

	router := BuildServiceResources()

	port := fmt.Sprintf(":%v", app.Config.ServerPort)
	log.Info("server started at port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func BuildServiceResources() *mux.Router {
	database := db.Postgres{}
	router := mux.NewRouter()

	dao := daos.NewLeiDAO(&database)
	leiService := service.NewLeiService(dao)

	api.ServeLeiResource(router, leiService)
	return router
}

func SetUpApplication() {
	log.Info("set up application")

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
