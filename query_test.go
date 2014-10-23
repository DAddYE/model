package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	q := Select("name", "surname").From("users").Where("id=?").Limit(1).Offset(0)
	assert.Equal(t, "SELECT name, surname FROM users WHERE id=? LIMIT 1 OFFSET 0", q.String())

	q = InsertInto("users", "name", "surname")
	assert.Equal(t, "INSERT INTO users ( name, surname ) VALUES ( ?, ? )", q.String())

	q = Update("users", "name=?", "surname=?")
	assert.Equal(t, "UPDATE users SET name=?, surname=?", q.String())
}
