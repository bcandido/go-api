package db

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"../app"
)

var DRIVER = "postgres"

var Connetion *sql.DB

func OpenConnection() *sql.DB {

	host := app.Config.DB["host"]
	port := app.Config.DB["port"]
	user := "postgres"
	password := "password"
	dbname := "postgres"

	fmt.Println(host, port)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)


	Connection, err := sql.Open(DRIVER, psqlInfo)
	if err != nil {
		panic(err)
	}

	return Connection
}