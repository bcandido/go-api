package daos

import (
	"github.com/op/go-logging"
	"../models"
	"../db"
	"errors"
	"fmt"
)

const MODULE = "daos"

var log = logging.MustGetLogger(MODULE)

var (
	ErrorNoItemFound   = errors.New("no item found")
	ErrorDataBaseConnection = errors.New("unable to establish a connection with the database")
	ErrorTransactionBegin = errors.New("could not begin a transaction")
	ErrorTransactionCommit = errors.New("could not commit insertion")
	ErrorLeiAlreadyInserted = errors.New("lei already insert")
)

// LeiDAO persists Lei data in database
type LeiDAO struct {
	database *db.Postgres
}

// NewLeiDAO creates a new LeiDAO
func NewLeiDAO(postgres *db.Postgres) *LeiDAO {
	return &LeiDAO{database: postgres}
}

// Get reads the Lei with the specified ID from the database.
func (dao *LeiDAO) GetAll() ([]models.Lei, error) {

	// open db connection
	err := dao.database.Open()
	defer dao.database.Close()
	if err != nil {
		log.Error(err.Error())
		return []models.Lei{}, ErrorDataBaseConnection
	}

	tx, err := dao.database.DB.Begin()
	if err != nil {
		log.Error(err.Error())
		return []models.Lei{}, ErrorTransactionBegin

	}
	query := "SELECT \"ID\", \"NOME\" FROM public.leis"

	rows, err := tx.Query(query)
	if err != nil {
		log.Error(err.Error())
		return []models.Lei{}, err
	}

	var leis []models.Lei
	for rows.Next() {
		var lei models.Lei
		rows.Scan(&lei.Id, &lei.Nome)
		leis = append(leis, lei)
	}
	return leis, nil
}


func (dao *LeiDAO) Get(id string) (models.Lei, error) {

	// open db connection
	err := dao.database.Open()
	defer dao.database.Close()
	if err != nil {
		log.Error(err.Error())
		return models.Lei{}, err
	}

	tx, err := dao.database.DB.Begin()
	if err != nil {
		log.Error(err.Error())
		return models.Lei{}, ErrorTransactionBegin

	}
	query := "SELECT \"ID\", \"NOME\" FROM public.leis WHERE \"ID\" ='" + id + "'"

	rows, err := tx.Query(query)
	if err != nil {
		log.Error(err.Error())
		return models.Lei{}, err
	}

	var lei models.Lei
	rows.Next()
	err = rows.Scan(&lei.Id, &lei.Nome)
	if err != nil {
		err = ErrorNoItemFound
	}
	return lei, err
}
func (dao *LeiDAO) Add(newLei string) error {

	// open db connection
	err := dao.database.Open()
	defer dao.database.Close()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	tx, err := dao.database.DB.Begin()
	if err != nil {
		log.Error()
		return ErrorTransactionBegin

	}
	query := fmt.Sprintf("INSERT INTO public.leis (\"NOME\") VALUES ('%s')", newLei)

	_, err = tx.Exec(query)
	if err != nil {
		log.Error(err.Error())
		if err = tx.Rollback(); err != nil {
			log.Critical(err.Error())
		}
		return ErrorLeiAlreadyInserted
	}

	if err = tx.Commit(); err != nil {
		log.Error(err.Error())
		if err = tx.Rollback(); err != nil {
			log.Critical(err.Error())
		}
		return ErrorTransactionCommit
	}

	return nil
}