package main

import (
	"fmt"
	"github.com/op/go-logging"
	"./app"
	"./db"
	//"./models"
	"net/http"
	"html"
	"github.com/gorilla/mux"
)

var log = logging.MustGetLogger("main")

func main() {
	log.Info("start application")

	// load application configurations
	log.Info("load application configurations")
	err := app.LoadConfig("./")
	if err != nil {
		log.Error("invalid application configuration")
		panic(err)
	}

	log.Info("open database connection")
	database := db.Postgres{}
	database.Open()

	defer log.Info("close database connection")
	defer database.Close()

	if err != nil {
		log.Error("not able to establish a connection with the database")
	}
	log.Info("connection establish")


	//rows, err := database.Query("SELECT id, nome FROM public.leis")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var leis []models.Lei
	//for rows.Next() {
	//	var lei models.Lei
	//	rows.Scan(&lei.Id, &lei.Nome)
	//	leis = append(leis, lei)
	//}
	//
	//for _, lei := range leis {
	//	fmt.Println(lei.Nome)
	//}


	port := fmt.Sprintf(":%v", app.Config.ServerPort)


	router := mux.NewRouter()
	router.HandleFunc("/", Index)

	log.Info("server " + app.Version + " is started at " + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}