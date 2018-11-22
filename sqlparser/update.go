package sqlparser

import (
	"github.com/alecthomas/participle"
	"github.com/xumc/miniDB/store"
)

type SetValue struct {
	String  *string  `@String`
	Number  *int64   `| @Number`
	Boolean *Boolean ` | @("TRUE" | "FALSE")`
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

func (p *Parser) ParseUpdate(sql string) (*UpdateSQL, error) {
	ast := &UpdateSQL{}

	err := updateSqlParser.ParseString(sql, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (p *Parser) TransformUpdate(ast *UpdateSQL, tableDesc *store.TableDesc) (*store.QueryTree, []store.SetItem) {
	qt := transformWhere(ast.Where)
	storeSetitems := make([]store.SetItem, len(ast.SetItems))
	for i, si := range ast.SetItems {
		storeSetitems[i] = store.SetItem{
			Name: *si.Key,
			Value: func(record store.Record) (interface{}, error) {
				// TODO support function
				return transformSetValue(*si.Value), nil
			},
		}

	}
	return qt, storeSetitems
}

func transformWhere(where *QueryTree) *store.QueryTree {
	qt := store.QueryTree{
		Negative: where.Negative,
		MatchAll: (bool)(where.MatchAll),
	}

	if where.Item != nil {
		qt.Item = &store.QueryItem{
			Key:      *where.Item.Key,
			Operator: transformOperator(*where.Item.Operator),
			Value:    transformQueryValue(*where.Item.Value),
		}

		return &qt
	}

	qt.Left = transformWhere(where.LeftTree)
	qt.Right = transformWhere(where.RightTree)

	return &qt
}

func transformOperator(operator Operator) store.Matcher {
	switch operator {
	case Operator("="):
		return store.MatcherEqual{}
	case Operator("<"):
		return store.MatcherLessThan{}
	}

	return nil
}

func transformQueryValue(v QueryValue) interface{} {
	if v.String != nil {
		return *v.String
	}

	if v.Number != nil {
		return *v.Number
	}

	if v.Boolean != nil {
		return (bool)(*v.Boolean)
	}

	return nil
}

func transformSetValue(v SetValue) interface{} {
	if v.String != nil {
		return *v.String
	}

	if v.Number != nil {
		return *v.Number
	}

	if v.Boolean != nil {
		return (bool)(*v.Boolean)
	}

	return nil
}
