package sqlparser

import (
	"github.com/alecthomas/participle"
	"github.com/xumc/miniDB/store"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

type InsertField struct {
	Name *string `@Ident [","]`
}

type InsertValue struct {
	String  *string  `@String [","]`
	Number  *int64   `| @Number [","]`
	Boolean *Boolean ` | @("TRUE" | "FALSE") [","]`
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

func (p parser) TransformInsert(ast *InsertSQL, tableDesc *store.TableDesc) store.Record {
	record := store.Record{
		TableName: *ast.TableName,
	}

	record.Values = make([]interface{}, len(tableDesc.Columns))
	for i, c := range tableDesc.Columns {
		for fi, f := range ast.Fields {
			if *f.Name == c.Name {
				var v interface{}
				switch c.Type {
				case store.ColumnTypeInteger:
					v = *ast.Values[fi].Number
				case store.ColumnTypeString:
					v = *ast.Values[fi].String
				case store.ColumnTypeBool:
					v = (bool)(*ast.Values[fi].Boolean)
				case store.ColumnTypeByte: // TODO
				}
				record.Values[i] = v
			}
		}
	}

	return record
}
