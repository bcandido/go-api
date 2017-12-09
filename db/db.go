package db

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"../app"
	"github.com/op/go-logging"
)

const DRIVER = "postgres"

var log = logging.MustGetLogger("cofing")

type Postgres struct {
	DB *sql.DB
}

func (db *Postgres) Open() error {
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

	var err error
	db.DB, err = sql.Open(DRIVER, psqlInfo)
	if err != nil {
		log.Error("failure to open a connection:", err)
	} else {
		log.Info("connection opened successfully")
	}

	return err
}

func (db *Postgres) Close() {
	log.Info("closing connection")

	err := db.DB.Close()
	if err != nil {
		log.Error("failure to close a connection:", err)
		panic(err)
	}
	log.Info("connection closed")
}