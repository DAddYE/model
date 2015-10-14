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
			Skipped bool   `custom:"-"`
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
	assert.True(t, reflect.ValueOf(m.Fields[0].Interface).Elem().CanAddr())
	assert.Equal(t, 9, m.Values()[0])
}

func TestSqlNull(t *testing.T) {
	f := &struct {
		Name    sql.NullString `sql:"name"`
		Surname string         `sql:"surname"`
		Noop    bool
		Login   struct {
			Email    string `sql:"login"`
			Password string `sql:"password"`
		}
	}{}
	f.Name.String = "Stan"
	f.Surname = "Smith"
	f.Login.Email = "smith@gmail.com"

	m := New(f, "sql")
	assert.Len(t, m.Interfaces(), 4)
	assert.Equal(t, "Stan", m.Map()["name"].(sql.NullString).String)
	assert.Equal(t, "smith@gmail.com", m.Map()["login"].(string))
}
