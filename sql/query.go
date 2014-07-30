package query

import (
	"database/sql"

	"github.com/daddye/model"
)

type Iter struct {
	r *sql.Rows
	m *model.Model
}

func (iter *Iter) Next() (res bool) {
	if res = iter.r.Next(); res {
		iter.m.Err = iter.r.Scan(iter.m.Values...)
	}
	return
}

func First(m *model.Model, conditions string, args ...interface{}) error {
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions, "LIMIT 1"}

	m.Err = m.Interface.(*sql.DB).QueryRow(query.String(), args...).Scan(m.Values...)
	return m.Err
}

func Find(m *model.Model, conditions string, args ...interface{}) (*Iter, error) {
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions}

	rows, err := m.Interface.(*sql.DB).Query(query.String(), args...)
	if err != nil {
		m.Err = err
		return nil, err
	}

	return &Iter{m: m, r: rows}, nil
}
