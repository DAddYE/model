package model

import (
	"fmt"
	"reflect"
	"strconv"
)

// Model is an utility type to access and manipulate struct informations.
type Model struct {
	Fields []Field // field fields (pointers to fields)
	ref    interface{}
	tag    string
}

func structType(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("the be a struct or a pointer to a struct; got %T", v.Type()))
	}

	return v
}

// Allocates a new Model and will extract and cache informations of fields that have the given `tag`
func New(s interface{}, tag string) Model {
	v := reflect.ValueOf(s)

	if v.Type().Kind() != reflect.Ptr || v.Type().Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("the model must be pointer to a struct; got %T", v.Type()))
	}

	m := Model{
		ref: s,
		tag: tag,
	}

	sv, st := fields(structType(s), tag)

	m.Fields = make([]Field, len(sv))

	for i := 0; i < len(sv); i++ {
		m.Fields[i] = Field{
			TagName:     st[i].Tag.Get(tag),
			StructField: st[i],
			Interface:   sv[i].Addr().Interface(),
		}
	}

	return m
}

// Returns the real values of the tagged fields. This cannot be cached.
func Values(s interface{}, tag string) []interface{} {
	sv, st := fields(structType(s), tag)
	values := make([]interface{}, 0, len(sv))
	for i, value := range sv {
		if value.Type().Kind() == reflect.Struct && st[i].Tag.Get(tag) == "" {
			values = append(values, Values(value.Interface(), tag)...)
			continue
		}
		values = append(values, value.Interface())
	}
	return values
}

// Returns a map of names (tags) and their values
//	Example: m.Map()["last_name"]
func (m Model) Map() map[string]interface{} {
	values := Values(m.ref, m.tag)
	ret := make(map[string]interface{}, len(values))
	for i, value := range values {
		ret[m.TagNames()[i]] = value
	}
	return ret
}

// Returns an array of values of tagged fields
func (m Model) Values() []interface{} {
	return Values(m.ref, m.tag)
}

// Returns the "real" name of all struct's fields
func (m Model) Names() []string {
	ret := make([]string, len(m.Fields))
	for i, field := range m.Fields {
		ret[i] = field.Name
	}
	return ret
}

// Returns the tagged names of all struct's fields
func (m Model) TagNames() []string {
	ret := make([]string, len(m.Fields))
	for i, field := range m.Fields {
		ret[i] = field.TagName
	}
	return ret
}

// Returns the given string as many times as the len of model.Fields
// Useful when building insert statements with Query
// Example:
//	people.Repeat("?")
//	["?", "?", "?"] // len is the number of mapped people struct's fields
func (m Model) Repeat(s string) []string {
	ret := make([]string, len(m.Fields))
	for i, _ := range m.Fields {
		ret[i] = s
	}
	return ret
}

// Returns the given string as many times as the len of model.Fields plus his increment
// Useful when building insert statements (like in postgres) with Query struct.
// Example:
//	people.Repeat("$")
//	["$1", "$2", "$3"]
func (m Model) RepeatInc(s string) []string {
	ret := make([]string, len(m.Fields))
	for i, _ := range m.Fields {
		ret[i] = s + strconv.Itoa(i+1)
	}
	return ret
}

// Returns an array of assignment strings:
// Example
//	["name=?" "surname=?"]
func (m Model) Assing() []string {
	ret := make([]string, len(m.Fields))
	for i, field := range m.Fields {
		// I'm not sure for this use case this is more performant:
		// string(append([]byte(field.TagName), "=?"...))
		ret[i] = field.TagName + "=?"
	}
	return ret
}

// Returns the underlining interface of each struct's field.
// Useful when binding results to our struct.
func (m Model) Interfaces() []interface{} {
	ret := make([]interface{}, len(m.Fields))
	for i, field := range m.Fields {
		ret[i] = field.Interface
	}
	return ret
}

// Returns a map of names (tags) and their interfaces
//	Example: m.MapInterface()["last_name"]
func (m Model) MapInterface() map[string]interface{} {
	interfaces := m.Interfaces()
	ret := make(map[string]interface{}, len(interfaces))
	for i, value := range interfaces {
		ret[m.TagNames()[i]] = value
	}
	return ret
}
