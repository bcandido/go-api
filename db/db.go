package db

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"../app"
	"github.com/op/go-logging"
)

const DRIVER = "postgres"
const MODULE = "config"

var log = logging.MustGetLogger(MODULE)

type Postgres struct {
	DB *sql.DB
}

type Tx struct {
	Tx *sql.Tx
}

type Query struct {
	Query  string
	Params []interface{}
}

func (db *Postgres) Open() error {

	host := app.Config.DB["host"]
	port := app.Config.DB["port"]
	user := "postgres"
	password := "password"
	dbname := "postgres"

	credentials := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	credentials = fmt.Sprintf(credentials, host, port, user, password, dbname)
	var err error
	db.DB, err = sql.Open(DRIVER, credentials)
	if err != nil {
		log.Error("failure to open a connection:", err)
	}

	return err
}

func (db *Postgres) Close() {
	err := db.DB.Close()
	if err != nil {
		log.Error("failure to close database connection:", err)
	}
}
