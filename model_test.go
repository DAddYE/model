package model

import (
	"reflect"
	"testing"
)

type Feed struct {
	Id         int    `sql:"id"`
	SourceType string `sql:"source_type"`
	Model
}

func TestModel(t *testing.T) {
	f := new(Feed)
	SetModel(f, "sql")

	cols := []string{"id", "source_type"}

	if l1, l2 := len(f.Columns), len(cols); l1 != l2 {
		t.Errorf("expected %d columns, got: %d", l2, l1)
	}

	for i, col := range cols {
		if f.Columns[i] != col {
			t.Errorf("expected column %d to be `%s`, intead got: `%s`", i+1, col, f.Columns[i])
		}
	}

	if l1, l2 := len(f.Values), len(cols); l1 != l2 {
		t.Errorf("expected %d values, got: %d", l2, l1)
	}

	if typ := reflect.TypeOf(f.Values[0]); typ.Kind() != reflect.Ptr ||
		typ.Elem().Kind() != reflect.Int {
		t.Errorf("expected value %d to be a `%s`, got: `%s`", 1, reflect.Int, typ.Elem().Kind())
	}

	if typ := reflect.TypeOf(f.Values[1]); typ.Kind() != reflect.Ptr ||
		typ.Elem().Kind() != reflect.String {
		t.Errorf("expected value %d to be a `%s`, got: `%s`", 2, reflect.Int, typ.Elem().Kind())
	}
}
