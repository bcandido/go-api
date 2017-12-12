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
	ErrorTransactionFailure = errors.New("unable to begin a transaction")
	ErrorInsertionFailure = errors.New("unable to get leis data")
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
		log.Error("unable to get leis data")
		return []models.Lei{}, ErrorTransactionFailure

	}
	query := "SELECT \"ID\", \"NOME\" FROM public.leis"

	rows, err := tx.Query(query)
	if err != nil {
		log.Error("unable to get leis data")
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
		log.Error("unable to get leis data")
		return models.Lei{}, ErrorTransactionFailure

	}
	query := "SELECT id, nome FROM public.leis WHERE id = '" + id + "'"

	rows, err := tx.Query(query)
	if err != nil {
		log.Error("unable to get leis data")
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
func (dao *LeiDAO) Add(newLei string) (bool, error) {

	// open db connection
	err := dao.database.Open()
	defer dao.database.Close()
	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	tx, err := dao.database.DB.Begin()
	if err != nil {
		log.Error()
		return false, ErrorTransactionFailure

	}
	query := fmt.Sprintf("INSERT INTO public.leis (\"NOME\") VALUES ('%s')", newLei)

	rows, err := tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return false, ErrorInsertionFailure
	}

	tx.Commit()

	id, _ := rows.LastInsertId()
	log.Info("LastInsertId:", id)

	return true, err
}