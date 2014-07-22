package query

import "github.com/daddye/model"

func First(m *model.Model, conditions string, args ...interface{}) error {
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions, "LIMIT 1"}

	m.Err = m.DB.QueryRow(query.String(), args...).Scan(m.Values...)
	return m.Err
}

func Find(m *model.Model, conditions string, args ...interface{}) (*model.Iter, error) {
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions}

	rows, err := m.DB.Query(query.String(), args...)
	if err != nil {
		m.Err = err
		return nil, err
	}

	iter := model.Iter{Model: m, Rows: rows}
	return &iter, nil
}
