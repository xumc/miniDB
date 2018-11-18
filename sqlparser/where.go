package sqlparser

import "fmt"

type MatchAll bool

func (m *MatchAll) Capture(values []string) error {
	if values[0] != "AND" && values[0] != "OR" {
		return fmt.Errorf("invalid identifier %s", values[0])
	}
	*m = values[0] == "AND"
	return nil
}

type Operator string

func (m *Operator) Capture(values []string) error {
	*m = (Operator)(values[0])
	return nil
}

type QueryTree struct {
	Negative  bool       `[@"!"]`
	LeftTree  *QueryTree `"(" @@`
	MatchAll  MatchAll   `@("AND"|"OR")`
	RightTree *QueryTree `@@ ")"`
	Item      *QueryItem `| @@`
}

type QueryValue struct {
	String  *string  `@String`
	Number  *int64   `| @Number`
	Boolean *Boolean ` | @("TRUE" | "FALSE")`
}

type QueryItem struct {
	Key      *string     `@Ident`
	Operator *Operator   `@("="|"<")`
	Value    *QueryValue `@@`
}
