package transaction

import (
	"log"
	"sync"

	"github.com/xumc/miniDB/store"
)

type Transactor interface {
	Lock()
	UnLock()
	Rollback()
}

type Transaction struct {
	logger *log.Logger
	store  *store.Store

	sync.Mutex // global lock
}

func NewTransaction(logger *log.Logger, s *store.Store) *Transaction {
	return &Transaction{
		logger: logger,
		store:  s,
	}
}

func (t *Transaction) RegisterTable(tableDesc store.TableDesc) error {
	// TODO
	return nil
}

func (t *Transaction) Insert(tableName string, record store.Record) (affectedRows int64, err error) {
	return t.store.Insert(tableName, record)
}

func (t *Transaction) Select(tableName string, qt *store.QueryTree) ([]store.Record, error) {
	return t.store.Select(tableName, qt)
}

func (t *Transaction) Update(tableName string, qt *store.QueryTree, setItems []store.SetItem) (affectedRows int64, err error) {
	return t.store.Update(tableName, qt, setItems)
}

func (t *Transaction) Delete(talbeName string, qt *store.QueryTree) (affectedRows int64, err error) {
	return t.Delete(talbeName, qt)
}

func (t *Transaction) Rollback() {

}

func (t *Transaction) Run() {

}
