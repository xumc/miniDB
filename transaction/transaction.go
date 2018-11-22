package transaction

import (
	"log"

	"github.com/xumc/miniDB/store"
)

type Transactor interface {
	Lock()
	UnLock()
	Rollback()

	Run()
}

type Transaction struct {
	logger *log.Logger
	store  *store.Store
}

func NewTransaction(logger *log.Logger, s *store.Store) *Transaction {
	return &Transaction{
		logger: logger,
		store:  s,
	}
}

func (s *Transaction) Lock() {
}

func (s *Transaction) UnLock() {
}

func (s *Transaction) Rollback() {

}

func (s *Transaction) Run() {

}
