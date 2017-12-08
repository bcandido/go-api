package main

import (
	"fmt"
	"github.com/op/go-logging"
	"./app"
	"./db"
	"./models"
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
	conn := db.OpenConnection()

	defer log.Info("close database connection")
	defer conn.Close()

	if err != nil {
		log.Error("not able to establish a connection with the database")
	}
	log.Info("connection establish")


	rows, err := conn.Query("SELECT id, nome FROM public.leis")
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


}
