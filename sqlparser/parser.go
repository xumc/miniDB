package sqlparser

import (
	"fmt"
	"log"
	"strings"

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
