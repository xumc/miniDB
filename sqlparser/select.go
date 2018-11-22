package sqlparser

import (
	"github.com/alecthomas/participle"
	"github.com/xumc/miniDB/store"
)

type FieldValue struct {
	String    *string `@Ident`
	AllFields bool    `| @"*"`
}

type SelectSQL struct {
	Type      *string       `@Ident`
	Fields    []*FieldValue `@@ { "," @@ }`
	TableName *string       `"FROM" @Ident`
	Where     *QueryTree    `"WHERE" @@ ";"`
}

var selectSqlParser = participle.MustBuild(
	&SelectSQL{},
	participle.Lexer(sqlLexer),
	participle.Unquote("String"),
	participle.CaseInsensitive("Ident"),
	participle.Elide("Whitespace", "Comment"),
)

func (p *Parser) ParseSelect(sql string) (*SelectSQL, error) {
	ast := &SelectSQL{}
	err := selectSqlParser.ParseString(sql, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (p *Parser) TransformSelect(ast *SelectSQL) *store.QueryTree {
	return transformWhere(ast.Where)
}
