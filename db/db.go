package db

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"../app"
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"syscall"
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

	// close the connection for SIGTERM
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.Close()
		os.Exit(1)
	}()

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
