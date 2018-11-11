package sqlparser

import (
	"github.com/xumc/miniDB/store"
)

type Parser interface {
	Parse(sql string) (tableName string, setItems []store.SetItem, qt *store.QueryTree, err error)
}

type parse struct{}
