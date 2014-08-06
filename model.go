package model

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Model struct {
	Fields    []*Field // field fields (pointers to fields)
	reference interface{}
	tag       string
}

type Field struct {
	TagName    string
	Reflection reflect.Value
	Interface  interface{}
	reflect.StructField
}

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
		default:
			panic(fmt.Errorf("the type %T is invalid", x))
		}
	}
	return strings.Join(res, " ")
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

func New(s interface{}, tag string) *Model {
	v := reflect.ValueOf(s)

	if v.Type().Kind() != reflect.Ptr || v.Type().Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("the model must be pointer to a struct; got %T", v.Type()))
	}

	m := new(Model)
	m.reference = s
	m.tag = tag

	sv, st := fields(structType(s))

	m.Fields = make([]*Field, 0, len(sv))

	for i := 0; i < len(sv); i++ {
		// check if is tagged
		name := st[i].Tag.Get(tag)
		if name == "" {
			continue
		}

		field := &Field{
			TagName:     st[i].Tag.Get(tag),
			StructField: st[i],
			Reflection:  sv[i],
			Interface:   sv[i].Addr().Interface(),
		}

		m.Fields = append(m.Fields, field)
	}

	return m
}

func fields(sv reflect.Value) ([]reflect.Value, []reflect.StructField) {
	v := make([]reflect.Value, 0)
	t := make([]reflect.StructField, 0)
	st := sv.Type()

	for i := 0; i < st.NumField(); i++ {
		// check if we are in a embedded struct
		if st.Field(i).Type.Kind() == reflect.Struct {
			vn, tn := fields(sv.Field(i))
			v = append(v, vn...)
			t = append(t, tn...)
			continue
		}
		v = append(v, sv.Field(i))
		t = append(t, st.Field(i))
	}

	if len(v) != len(t) {
		panic("internal error")
	}

	return v, t
}

func Values(s interface{}, tag string) []interface{} {
	sv, st := fields(structType(s))
	values := make([]interface{}, 0, len(sv))
	for i, value := range sv {
		if st[i].Tag.Get(tag) == "" {
			continue
		}
		if value.Type().Kind() == reflect.Struct {
			values = append(values, Values(value.Interface(), tag)...)
			continue
		}
		values = append(values, value.Interface())
	}
	return values
}

func (m *Model) Map() map[string]interface{} {
	values := Values(m.reference, m.tag)
	ret := make(map[string]interface{}, len(values))
	for i, value := range values {
		ret[m.TagNames()[i]] = value
	}
	return ret
}

func (m *Model) Values() []interface{} {
	return Values(m.reference, m.tag)
}

func (m *Model) Names() []string {
	ret := make([]string, 0, len(m.Fields))
	for _, field := range m.Fields {
		ret = append(ret, field.Name)
	}
	return ret
}

func (m *Model) TagNames() []string {
	ret := make([]string, 0, len(m.Fields))
	for _, field := range m.Fields {
		ret = append(ret, field.TagName)
	}
	return ret
}

func (m *Model) Interfaces() []interface{} {
	ret := make([]interface{}, 0, len(m.Fields))
	for _, field := range m.Fields {
		ret = append(ret, field.Interface)
	}
	return ret
}

func (m *Model) Reflections() []reflect.Value {
	ret := make([]reflect.Value, 0, len(m.Fields))
	for _, field := range m.Fields {
		ret = append(ret, field.Reflection)
	}
	return ret
}

// alpha
func (m *Model) Decode(raw map[string]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	for k, v := range raw {
		for _, field := range m.Fields {
			if field.TagName == k {
				field.Reflection.Set(reflect.ValueOf(v))
			}
		}
	}
	return
}
