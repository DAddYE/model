package model

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	f := &struct {
		Id  int `custom:"id"`
		Bio struct {
			Name    string `custom:"name"`
			Surname string `custom:"surname"`
		}
		Login   string
		private bool
	}{}
	f.Id = 9
	f.Bio.Name = "Stan"
	f.Bio.Surname = "Smith"

	// attach our model
	m := New(f, "custom")

	// we expect 3 fields, because they are the only one mapped
	assert.Equal(t, 3, len(m.Fields))

	assert.Equal(t, "id", m.Fields[0].TagName)
	assert.Equal(t, reflect.Int, m.Fields[0].Reflection.Kind())
	assert.True(t, reflect.ValueOf(m.Fields[0].Interface).Elem().CanAddr())
	assert.Equal(t, 9, m.Values()[0])
}

func TestDecode(t *testing.T) {
	f := &struct {
		Id  int `my:"id"`
		Bio struct {
			Name    string `my:"name"`
			Surname string `my:"last_name"`
		}
		private bool
	}{}
	f.Id = 9
	f.Bio.Name = "Stan"
	f.Bio.Surname = "Smith"
	m := New(f, "my")

	// try best scenario
	err := m.Decode(map[string]interface{}{"id": 18, "name": "Francine"})
	assert.NoError(t, err)
	assert.Equal(t, 18, f.Id)
	assert.Equal(t, "Francine", f.Bio.Name)
	assert.Equal(t, "Smith", f.Bio.Surname)

	// TODO: nullify, invalid conversion, etc...
}

func TestSqlNull(t *testing.T) {
	f := &struct {
		Name    sql.NullString `sql:"name"`
		Surname string         `sql:"surname"`
	}{}
	f.Name.String = "Stan"
	f.Surname = "Smith"

	m := New(f, "sql")
	assert.Len(t, m.Interfaces(), 2)
	assert.Equal(t, "Stan", m.Map()["name"].(sql.NullString).String)
}
