package db

import (
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"database/sql"
	"fmt"
	"../app"
	"context"
	"errors"
)

const DRIVER = "postgres"
const MODULE = "config"

const maxConnections = 30
const maxIdleConnections = 30

var ErrorConnection = errors.New("cloud not create a data base connection")

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

	dsn := DataSourceName{
		host:     fmt.Sprint(app.Config.DB["host"]),
		port:     fmt.Sprint(app.Config.DB["port"]),
		user:     "postgres",
		password: "password",
		dbName:   "postgres",
	}

	// open database connection
	var err error
	db.DB, err = sql.Open(DRIVER, dsn.GetDSN())
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// validate database connection
	if err = db.DB.Ping(); err != nil {
		log.Error(err.Error())
		panic(ErrorConnection)
	}

	// set connection pool idle/max connection
	db.DB.SetMaxOpenConns(maxConnections)
	db.DB.SetMaxIdleConns(maxIdleConnections)

	return nil
}

func (db *Postgres) Close() {
	err := db.DB.Close()
	if err != nil {
		log.Error("failure to close database connection:", err)
	}
}

func (db *Postgres) GetConnection() (*sql.Conn, error) {
	ctx := context.Background()
	conn, err := db.DB.Conn(ctx)
	return conn, err
}

func (db *Postgres) Select(query string) (*sql.Rows, error) {
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Error(err.Error())
		return &sql.Rows{}, err
	}

	return rows, nil
}

func (db *Postgres) Exec(ctx context.Context, query string) (*sql.Rows, error) {
	conn, err := db.GetConnection()
	defer conn.Close()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		log.Error(err.Error())
		return &sql.Rows{}, err

	}

	rows, err := tx.Query(query)
	if err != nil {
		log.Error(err.Error())
		if err = tx.Rollback(); err != nil {
			log.Critical(err.Error())
		}
		return &sql.Rows{}, err
	}

	if err = tx.Commit(); err != nil {
		log.Error(err.Error())
		if err = tx.Rollback(); err != nil {
			log.Critical(err.Error())
		}
		return &sql.Rows{}, err
	}

	return rows, nil
}

type DataSourceName struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func (dsn *DataSourceName) GetDSN() string {
	postgresDSNFormat := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	dataSourceName := fmt.Sprintf(postgresDSNFormat, dsn.host, dsn.port, dsn.user, dsn.password, dsn.dbName)
	log.Info(dataSourceName)
	return dataSourceName
}
