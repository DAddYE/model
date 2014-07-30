package query

import (
	"github.com/daddye/model"
	"github.com/gocql/gocql"
)

type Iter struct {
	i *gocql.Iter
	m *model.Model
}

func (iter *Iter) Next() (res bool) {
	return iter.i.Scan(iter.m.Values...)
}

func First(m *model.Model, conditions string, args ...interface{}) error {
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions, "LIMIT 1"}

	m.Err = m.Interface.(*gocql.Session).Query(query.String(), args...).Scan(m.Values...)
	return m.Err
}

func Find(m *model.Model, conditions string, args ...interface{}) *Iter {
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions}

	iter := m.Interface.(*gocql.Session).Query(query.String(), args...).Iter()
	return &Iter{m: m, i: iter}
}
