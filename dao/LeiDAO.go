package dao

import (
	"database/sql"
	"github.com/op/go-logging"
	"../models"
)

const MODULE = "dao"

var log = logging.MustGetLogger(MODULE)

// LeiDAO persists Lei data in database
type LeiDAO struct {
	database *sql.DB
}

// NewLeiDAO creates a new LeiDAO
func NewLeiDAO(database *sql.DB) *LeiDAO {
	return &LeiDAO{database: database}
}

// Get reads the Lei with the specified ID from the database.
func (dao *LeiDAO) Get(id int) (error) {
	tx, _ := dao.database.Begin()

	query := "SELECT id, nome FROM public.leis WHERE id = " + string(id)
	leis, err := tx.Exec(query)
	if err != nil {
		log.Error("unable to get data")
		return err
	}

	log.Info(leis)

	return err
}

func (dao *LeiDAO) GetAll() ([]models.Lei, error) {
	tx, _ := dao.database.Begin()

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

	log.Info(leis)

	return leis, err
}
