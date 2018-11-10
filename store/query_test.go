package store

import (
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	// find students who didn't pass the exam, age < 18 and female
	qt := QueryTree{
		MatchAll: true,
		Negative: false,
		Left: &QueryTree{
			MatchAll: true,
			Negative: false,

			Left: &QueryTree{
				Item: &QueryItem{Key: "pass", Operator: MatcherEqual{}, Value: false},
			},
			Right: &QueryTree{
				Item: &QueryItem{Key: "age", Operator: MatcherLessThan{}, Value: int64(18)},
			},
		},
		Right: &QueryTree{
			Item: &QueryItem{Key: "sex", Operator: MatcherEqual{}, Value: "FEMALE"},
		},
	}
	if qt.PrettyPrint() != "((pass=false AND age<18 ) AND sex='FEMALE' ) " {
		t.Fail()
	}
}

func TestIsQueryTreeMatch(t *testing.T) {
	tableDesc := &TableDesc{
		Name: "student",
		Columns: []Column{
			Column{Name: "id", Type: ColumnTypeInteger},
			Column{Name: "name", Type: ColumnTypeString},
			Column{Name: "age", Type: ColumnTypeInteger},
			Column{Name: "pass", Type: ColumnTypeBool},
		},
	}

	t.Run("case 1, simple query", func(t *testing.T) {
		qt := &QueryTree{
			Item: &QueryItem{Key: "name", Operator: MatcherEqual{}, Value: "jack"},
		}
		recordValue := []interface{}{int64(1), "jack", int64(10), false}
		match, err := isQueryTreeMatch(tableDesc, qt, recordValue)
		if err != nil {
			t.Fatalf("case 1 failed due to error")
		}
		if !match {
			t.Fatalf("case 1 failed")
		}
	})

	t.Run("case 2, empty left tree", func(t *testing.T) {
		qt := &QueryTree{
			Right: &QueryTree{
				Item: &QueryItem{Key: "name", Operator: MatcherEqual{}, Value: "jack"},
			},
		}
		recordValue := []interface{}{int64(1), "jack", int64(10), false}
		match, err := isQueryTreeMatch(tableDesc, qt, recordValue)
		if err != nil {
			t.Fatalf("case 2 failed due to error")
		}
		if !match {
			t.Fatalf("case 2 failed")
		}
	})

	t.Run("case 3, negative case", func(t *testing.T) {
		qt := &QueryTree{
			Negative: true,
			Item:     &QueryItem{Key: "name", Operator: MatcherEqual{}, Value: "jack"},
		}
		recordValue := []interface{}{int64(1), "jack", int64(10), false}
		match, err := isQueryTreeMatch(tableDesc, qt, recordValue)
		if err != nil {
			t.Fatalf("case 3 failed due to error")
		}
		if match {
			t.Fatalf("case 3 failed")
		}
	})

	t.Run("case 4, left and right are all simple tree", func(t *testing.T) {
		qt := &QueryTree{
			Left: &QueryTree{
				Item: &QueryItem{Key: "name", Operator: MatcherEqual{}, Value: "jack"},
			},
			Right: &QueryTree{
				Item: &QueryItem{Key: "age", Operator: MatcherLessThan{}, Value: int64(18)},
			},
		}
		recordValue := []interface{}{int64(1), "jack", int64(10), false}
		match, err := isQueryTreeMatch(tableDesc, qt, recordValue)
		if err != nil {
			t.Fatalf("case 4 failed due to error")
		}
		if !match {
			t.Fatalf("case 4 failed")
		}
	})
}
