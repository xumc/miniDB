package sqlparser

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/xumc/miniDB/transaction"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
	"github.com/xumc/miniDB/store"
)

type SQLParser interface {
	Parse(sql string) (SQL, error)

	TransformInsert(ast *InsertSQL, tableDesc *store.TableDesc) store.Record
	TransformUpdate(ast *UpdateSQL, tableDesc *store.TableDesc) (*store.QueryTree, []store.SetItem)
	TransformSelect(ast *SelectSQL) *store.QueryTree
	TransformDelete(ast *DeleteSQL) *store.QueryTree

	Run()
}

type Parser struct {
	logger      *log.Logger
	transaction *transaction.Transaction
}

func NewParser(logger *log.Logger, t *transaction.Transaction) *Parser {
	return &Parser{
		logger:      logger,
		transaction: t,
	}
}

func (p *Parser) Run() {

}

type SQL interface{}

var sqlLexer = lexer.Must(ebnf.New(`
Comment = "--" { "\u0000"…"\uffff"-"\n" } .
Ident = (alpha | "_") { "_" | alpha | digit } .
String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
Number = [ "-" | "+" ] ("." | digit) {"." | digit} .
Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
Whitespace = " " | "\t" | "\n" | "\r" .
alpha = "a"…"z" | "A"…"Z" .
digit = "0"…"9" .
any = "\u0000"…"\uffff" .
	`))

func (p *Parser) Parse(sql string) (SQL, error) {
	trimSQL := strings.TrimSpace(sql)

	if strings.HasPrefix(trimSQL, "INSERT") {
		return p.ParseInsert(trimSQL)
	} else if strings.HasPrefix(trimSQL, "SELECT") {
		return p.ParseSelect(trimSQL)
	} else if strings.HasPrefix(trimSQL, "UPDATE") {
		return p.ParseUpdate(trimSQL)
	} else if strings.HasPrefix(trimSQL, "DELETE") {
		return p.ParseDelete(trimSQL)
	} else {
		return nil, fmt.Errorf("unsupported sql %s", trimSQL)
	}
}

type ColumnHeader struct {
	Name string
	Type store.ColumnTypes
}

type Result interface {
	GetAffectedRows() int
	GetResultData() [][]interface{}
	GetResultHeader() []ColumnHeader
}

type TableResult struct {
	data []store.Record
}

func (tr *TableResult) GetAffectedRows() int {
	return int(len(tr.data))
}

func (tr *TableResult) GetResultData() [][]interface{} {
	ret := make([][]interface{}, len(tr.data))
	for i, c := range tr.data {
		ret[i] = c.Values
	}
	return ret
}

func (tr *TableResult) GetResultHeader() []ColumnHeader {
	return []ColumnHeader{} // TODO
}

type AffectedRowResult struct {
	affectedRows int
}

func (tr *AffectedRowResult) GetAffectedRows() int {
	return tr.affectedRows
}

func (tr *AffectedRowResult) GetResultData() [][]interface{} {
	return nil
}

func (tr *AffectedRowResult) GetResultHeader() []ColumnHeader {
	return nil
}

type beginSQL struct{}
type commitSQL struct{}

var BeginSQL = beginSQL{}
var CommitSQL = commitSQL{}

func (p *Parser) Next(tid uuid.UUID, sql SQL) (Result, error) {
	switch sqlStruct := sql.(type) {
	case beginSQL:
		p.transaction.Lock()
		return nil, nil
	case commitSQL:
		p.transaction.Unlock()
		return nil, nil
	case *InsertSQL:
		tableDesc, err := store.GetMetadataOf(*sqlStruct.TableName)
		if err != nil {
			return nil, err
		}
		record := p.TransformInsert(sqlStruct, tableDesc)

		affectedRows, err := p.transaction.Insert(record.TableName, record)
		if err != nil {
			return nil, err
		}

		p.logger.Printf("affected rows: %d", affectedRows)

	case *UpdateSQL:
		tableDesc, err := store.GetMetadataOf(*sqlStruct.TableName)
		if err != nil {
			return nil, err
		}

		qt, setItems := p.TransformUpdate(sqlStruct, tableDesc)

		affectedRows, err := p.transaction.Update(*sqlStruct.TableName, qt, setItems)
		if err != nil {
			return nil, err
		}

		p.logger.Printf("affected rows: %d", affectedRows)

	case *SelectSQL:
		qt := p.TransformSelect(sqlStruct)

		rs, err := p.transaction.Select(*sqlStruct.TableName, qt)
		if err != nil {
			return nil, err
		}

		p.logger.Printf("rows: %v", rs)

	case *DeleteSQL:
		qt := p.TransformDelete(sqlStruct)

		affectedRows, err := p.transaction.Delete(*sqlStruct.TableName, qt)
		if err != nil {
			return nil, err
		}

		p.logger.Printf("affected rows: %d", affectedRows)

	default:
		return nil, errors.New("unsupport sql type")
	}

	return nil, nil
}
