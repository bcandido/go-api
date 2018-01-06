package lei

import (
	"errors"
	"fmt"
	"context"
	"../db"
)

type Lei struct {
	Id   int    `json:"id" db:"id"`
	Nome string `json:"nome" db:"nome"`
}


var (
	ErrorNoRowsFound        = errors.New("no rows found")
	ErrorTransactionBegin   = errors.New("could not begin a transaction")
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
func (dao *LeiDAO) GetAll() ([]Lei, error) {

	query := "SELECT \"ID\", \"NOME\" FROM public.leis"
	rows, err := dao.database.Select(query)
	if err != nil {
		return []Lei{}, err
	}

	var leis []Lei
	for rows.Next() {
		var lei Lei
		rows.Scan(&lei.Id, &lei.Nome)
		leis = append(leis, lei)
	}
	return leis, nil
}

func (dao *LeiDAO) Get(id string) (Lei, error) {

	format := "SELECT \"ID\", \"NOME\" FROM public.leis WHERE \"ID\" = '%s'"
	query := fmt.Sprintf(format, id)
	log.Info(query)
	rows, err := dao.database.Select(query)
	if err != nil {
		return Lei{}, err
	}

	var lei Lei
	rows.Next()
	err = rows.Scan(&lei.Id, &lei.Nome)
	if err != nil {
		log.Info(err.Error())
		return Lei{}, ErrorNoRowsFound
	}
	return lei, nil
}
func (dao *LeiDAO) Add(newLei string) error {

	// create context
	ctx := context.Background()
	defer ctx.Done()

	format := "INSERT INTO public.leis (\"NOME\") VALUES ('%s')"
	query := fmt.Sprintf(format, newLei)

	_, err := dao.database.Exec(ctx, query)
	if err != nil {
		log.Error(err.Error())
		return ErrorLeiAlreadyInserted
	}

	return nil
}
