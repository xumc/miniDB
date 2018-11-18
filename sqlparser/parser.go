package sqlparser

import (
	"log"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
	"github.com/xumc/miniDB/store"
)

type Parser interface {
	Parse(sql string) (SQL, error)

	TransformInsert(ast *InsertSQL, tableDesc *store.TableDesc) store.Record
	TransformUpdate(ast *UpdateSQL, tableDesc *store.TableDesc) (*store.QueryTree, []store.SetItem)
	TransformSelect(ast *SelectSQL) *store.QueryTree
	TransformDelete(ast *DeleteSQL) *store.QueryTree
}

func NewParser(logger *log.Logger) Parser {
	return parser{logger: logger}
}

type SQL interface{}

type parser struct {
	logger *log.Logger
}

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

func (p parser) Parse(sql string) (SQL, error) {
	trimSQL := strings.TrimSpace(sql)

	if strings.HasPrefix(trimSQL, "INSERT") {
		return p.ParseInsert(trimSQL)
	} else if strings.HasPrefix(trimSQL, "SELECT") {
		return p.ParseSelect(trimSQL)
	} else if strings.HasPrefix(trimSQL, "UPDATE") {
		return p.ParseUpdate(trimSQL)
	} else if strings.HasPrefix(trimSQL, "DELETE") {
		return p.ParseDelete(trimSQL)
	}

	return nil, nil
}
