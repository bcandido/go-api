package service

import (
	"github.com/op/go-logging"
	"../models"
	"../daos"
)

const MODULE = "service"

var log = logging.MustGetLogger(MODULE)

type LeiService struct {
	dao *daos.LeiDAO
}

func NewLeiService(dao *daos.LeiDAO) *LeiService {
	return &LeiService{dao}
}

func (s *LeiService) GetAll() ([]models.Lei, error) {
	leis, err := s.dao.GetAll()
	if err != nil {
		log.Error("unable to get leis data from dao")
		return []models.Lei{}, err
	}
	return leis, nil
}

func (s *LeiService) Get(id string) (models.Lei, error) {
	lei, err := s.dao.Get(id)
	if err != nil {
		log.Error("unable to get leis data from dao")
		return models.Lei{}, err
	}
	return lei, nil
}

func (s *LeiService) Add(newLei string) (bool, error) {
	result, err := s.dao.Add(newLei)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	return result, nil
}
