package sqlparser

import (
	"github.com/alecthomas/participle"
	"github.com/xumc/miniDB/store"
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

func (p *Parser) ParseDelete(sql string) (*DeleteSQL, error) {
	ast := &DeleteSQL{}

	err := deleteSqlParser.ParseString(sql, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (p *Parser) TransformDelete(ast *DeleteSQL) *store.QueryTree {
	return transformWhere(ast.Where)
}
