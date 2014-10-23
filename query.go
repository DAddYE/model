package model

import (
	"fmt"
	"strconv"
	"strings"
)

// Query is a tiny helper to help you to build query strings.
// It will take care to join array of strings with ", " or convert `int`.
type Query []interface{}

func (q Query) String() string {
	res := make([]string, len(q))
	for i, v := range q {
		switch x := v.(type) {
		case string:
			res[i] = x
		case int:
			res[i] = strconv.Itoa(x)
		case []string:
			res[i] = strings.Join(x, ", ")
		case Query:
			res[i] = x.String()
		default:
			panic(fmt.Errorf("the type %T is invalid", x))
		}
	}
	return strings.Join(res, " ")
}

// in this way you can create your own version
var Placeholder = func(count int) []string {
	ret := make([]string, count)
	for i := 0; i < count; i++ {
		ret[i] = "?"
	}
	return ret
}

func Select(columns ...string) Query               { return Query{"SELECT", columns} }
func From(tables ...string) Query                  { return Query{"FROM", tables} }
func Where(conditions ...string) Query             { return Query{"WHERE", conditions} }
func Update(table string, columns ...string) Query { return Query{"UPDATE", table, "SET", columns} }
func Limit(value int) Query                        { return Query{"LIMIT", value} }
func Offset(value int) Query                       { return Query{"OFFSET", value} }
func InsertInto(table string, columns ...string) Query {
	return Query{"INSERT INTO", table, "(", columns, ") VALUES (", Placeholder(len(columns)), ")"}
}

func (q Query) Select(columns ...string) Query   { return append(q, Select(columns...)) }
func (q Query) From(tables ...string) Query      { return append(q, From(tables...)) }
func (q Query) Where(conditions ...string) Query { return append(q, Where(conditions...)) }
func (q Query) Limit(value int) Query            { return append(q, Limit(value)) }
func (q Query) Offset(value int) Query           { return append(q, Offset(value)) }
func (q Query) Update(table string, columns ...string) Query {
	return append(q, Update(table, columns...))
}
func (q Query) InsertInto(table string, columns ...string) Query {
	return append(q, InsertInto(table, columns...))
}
