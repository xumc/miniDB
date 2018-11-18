package sqlparser

import (
	"github.com/alecthomas/participle"
)

type DeleteSQL struct {
	Type      *string    `@Ident`
	TableName *string    `"FROM" @Ident`
	Where     *QueryTree `"WHERE" @@ ";"`
}

var deleteSqlParser = participle.MustBuild(
	&DeleteSQL{},
	participle.Lexer(sqlLexer),
	participle.Unquote("String"),
	participle.CaseInsensitive("Ident"),
	participle.Elide("Whitespace", "Comment"),
)

func (p parser) ParseDelete(sql string) (*DeleteSQL, error) {
	ast := &DeleteSQL{}

	err := deleteSqlParser.ParseString(sql, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}
