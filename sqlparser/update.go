package sqlparser

import (
	"github.com/alecthomas/participle"
)

type SetValue struct {
	String *string `@String`
	Number *int64  `| @Number`
}

type SetItem struct {
	Key   *string   `@Ident "="`
	Value *SetValue `@@`
}

type UpdateSQL struct {
	Type      *string    `@Ident`
	TableName *string    `@Ident`
	SetItems  []*SetItem `"SET" @@ { "," @@ }`
	// Where  []*QueryItem `"WHERE" {@@} ";"`
	Where *QueryTree `"WHERE" @@ ";"`
}

var updateSqlParser = participle.MustBuild(
	&UpdateSQL{},
	participle.Lexer(sqlLexer),
	participle.Unquote("String"),
	participle.CaseInsensitive("Ident"),
	participle.Elide("Whitespace", "Comment"),
)

func (p parser) ParseUpdate(sql string) (*UpdateSQL, error) {
	ast := &UpdateSQL{}

	err := updateSqlParser.ParseString(sql, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}
