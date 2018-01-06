package lei

import (
	"errors"
)

var (
	ErrorOperationFailure = errors.New("server fail to complete operation")
	ErrorNoItemFound      = errors.New("no item found")
	ErrorAlreadyInserted  = errors.New("already inserted")
)

type LeiService struct {
	dao *LeiDAO
}

func NewLeiService(dao *LeiDAO) *LeiService {
	return &LeiService{dao}
}

func (s *LeiService) GetAll() ([]Lei, error) {
	leis, err := s.dao.GetAll()
	if err != nil {
		log.Error(err.Error())
		return []Lei{}, err
	}
	return leis, nil
}

func (s *LeiService) Get(id string) (Lei, error) {
	lei, err := s.dao.Get(id)
	if err != nil {
		log.Error(err.Error())
		return Lei{}, err
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
	case ErrorNoRowsFound:
		return ErrorNoItemFound
	case ErrorLeiAlreadyInserted:
		return ErrorAlreadyInserted
	case ErrorTransactionBegin:
		return ErrorOperationFailure
	default:
		return err
	}
}
