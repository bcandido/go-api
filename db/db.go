package db

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"../app"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cofing")

const DRIVER = "postgres"

var Connetion *sql.DB

func OpenConnection() *sql.DB {
	log.Info("starting to open the connection")

	log.Info("loading configuration")
	host := app.Config.DB["host"]
	port := app.Config.DB["port"]
	user := "postgres"
	password := "password"
	dbname := "postgres"


	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Info(psqlInfo)

	Connection, err := sql.Open(DRIVER, psqlInfo)
	if err != nil {
		log.Error("failure to open a connection:", err)
		panic(err)
	}
	log.Info("connection opened successfully")

	return Connection
}