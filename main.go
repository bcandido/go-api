package main

import (
	"fmt"
	"github.com/op/go-logging"
	"./app"
	"./db"
	"./models"
	"net/http"
	"html"
	"github.com/gorilla/mux"
	"os"
	"os/signal"
	"syscall"
	"io/ioutil"
)

var log = logging.MustGetLogger("main")

func main() {
	log.Info("start application")
	log.Info("Version: " + app.Version)

	log.Info("set up application")
	SetUpApplication()

	log.Info("open database connection")
	database := db.Postgres{}
	err := database.Open()

	defer log.Info("close database connection")
	defer database.Close()

	if err != nil {
		log.Error("unable to establish a connection with the database")
	}
	log.Info("connection establish")


	rows, err := database.DB.Query("SELECT id, nome FROM public.leis")
	if err != nil {
		log.Fatal(err)
	}

	var leis []models.Lei
	for rows.Next() {
		var lei models.Lei
		rows.Scan(&lei.Id, &lei.Nome)
		leis = append(leis, lei)
	}

	for _, lei := range leis {
		fmt.Println(lei.Nome)
	}

	//
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		database.Close()
		os.Exit(1)
	}()

	router := mux.NewRouter()
	router.HandleFunc("/", Index)

	port := fmt.Sprintf(":%v", app.Config.ServerPort)
	log.Info("server version " + app.Version + " is started at " + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Proto, r.Host, r.Method, r.RequestURI)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func SetUpApplication() {
	// simply prints an ascii art
	PrintAsciiArt()

	// load application configurations
	LoadConfiguration()
}

func PrintAsciiArt()  {
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