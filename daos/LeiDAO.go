package daos

import (
	"github.com/op/go-logging"
	"../models"
	"../db"
)

const MODULE = "daos"

var log = logging.MustGetLogger(MODULE)

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
		message := "unable to establish a connection with the database"
		log.Error(message)
		return []models.Lei{}, err
	}

	tx, _ := dao.database.DB.Begin()
	query := "SELECT id, nome FROM public.leis"

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
	return leis, err
}


func (dao *LeiDAO) Get(id string) (models.Lei, error) {

	// open db connection
	err := dao.database.Open()
	defer dao.database.Close()
	if err != nil {
		message := "unable to establish a connection with the database"
		log.Error(message)
		return models.Lei{}, err
	}

	tx, _ := dao.database.DB.Begin()
	query := "SELECT id, nome FROM public.leis WHERE id = '" + id + "'"

	rows, err := tx.Query(query)
	if err != nil {
		log.Error("unable to get leis data")
		return models.Lei{}, err
	}

	var lei models.Lei
	err = rows.Scan(&lei.Id, &lei.Nome)
	if err != nil {
		log.Info("data not found for id = " + id)
	}
	return lei, err
}