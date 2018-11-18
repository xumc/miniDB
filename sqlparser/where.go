package sqlparser

import (
	"fmt"
	"strings"
)

type MatchAll bool

func (m *MatchAll) Capture(values []string) error {
	v := strings.ToLower(values[0])
	if v != "and" && v != "or" {
		return fmt.Errorf("invalid identifier %s", values[0])
	}
	*m = v == "and"
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
