package service

import (
	"github.com/op/go-logging"
	"../models"
	"../daos"
	"errors"
)

const MODULE = "service"

var log = logging.MustGetLogger(MODULE)

var (
	ErrorOperationFailure = errors.New("server fail to complete operation")
	ErrorNoItemFound      = errors.New("no item found")
	ErrorAlreadyInserted = errors.New("already inserted")
)

type LeiService struct {
	dao *daos.LeiDAO
}

func NewLeiService(dao *daos.LeiDAO) *LeiService {
	return &LeiService{dao}
}

func (s *LeiService) GetAll() ([]models.Lei, error) {
	leis, err := s.dao.GetAll()
	if err != nil {
		log.Error(err.Error())
		return []models.Lei{}, err
	}
	return leis, nil
}

func (s *LeiService) Get(id string) (models.Lei, error) {
	lei, err := s.dao.Get(id)
	if err != nil {
		log.Error(err.Error())
		return models.Lei{}, err
	}
	return lei, nil
}

func (s *LeiService) Add(newLei string) error {
	err := s.dao.Add(newLei)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func Validate(err error) error {
	switch err {
	case nil:
		return nil
	case daos.ErrorNoRowsFound:
		return ErrorNoItemFound
	case daos.ErrorLeiAlreadyInserted:
		return ErrorAlreadyInserted
	case daos.ErrorTransactionBegin:
			return ErrorOperationFailure
	default:
		return err
	}
}
