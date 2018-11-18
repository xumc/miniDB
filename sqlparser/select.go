package sqlparser

import (
	"github.com/alecthomas/participle"
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

func (p parser) ParseSelect(sql string) (*SelectSQL, error) {
	ast := &SelectSQL{}
	err := selectSqlParser.ParseString(sql, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}
