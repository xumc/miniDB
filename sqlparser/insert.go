package sqlparser

import (
	"github.com/alecthomas/participle"
)

type InsertField struct {
	Name *string `@Ident [","]`
}

type InsertValue struct {
	String *string `@String [","]`
	Number *int64  `| @Number [","]`
}

type InsertSQL struct {
	Type      *string        `@Ident`
	TableName *string        `"INTO" @Ident "("`
	Fields    []*InsertField `{@@} ")"`
	Values    []*InsertValue `"VALUES" "(" {@@} ")" ";"`
}

var insertSqlParser = participle.MustBuild(
	&InsertSQL{},
	participle.Lexer(sqlLexer),
	participle.Unquote("String"),
	participle.CaseInsensitive("Ident"),
	participle.Elide("Whitespace", "Comment"),
)

func (p parser) ParseInsert(sql string) (*InsertSQL, error) {
	ast := &InsertSQL{}
	err := insertSqlParser.ParseString(sql, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}
