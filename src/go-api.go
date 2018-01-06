package main

import (
	"fmt"
	"github.com/op/go-logging"
	"./lei"
	"./db"
	"./config"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
)

const MODULE = "main"

var log = logging.MustGetLogger(MODULE)

func main() {
	log.Info("start application...")
	SetUpApplication()

	router := BuildServiceResources()

	port := fmt.Sprintf(":%v", config.Config.ServerPort)
	log.Info("server started at port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func BuildServiceResources() *mux.Router {
	router := mux.NewRouter()

	dsn := db.NewDataSourceName(
		fmt.Sprint(config.Config.DB["host"]),
		fmt.Sprint(config.Config.DB["port"]),
		"postgres",
		"password",
		"postgres",
	)

	database := db.Postgres{}
	database.Open(dsn)

	dao := lei.NewLeiDAO(&database)
	leiService := lei.NewLeiService(dao)

	lei.ServeLeiResource(router, leiService)
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
	err := config.LoadConfig("./")
	if err != nil {
		log.Error("invalid application configuration")
		panic(err)
	}
	log.Info(config.Config)
}
