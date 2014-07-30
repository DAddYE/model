package model

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Columns []string
type Values []interface{}
type Query []interface{}

type Model struct {
	Columns
	Values
	Table     string
	Err       error
	Interface interface{} // a DB, gcql.Session... etc...
}

func (c Columns) String() string {
	return strings.Join(c, ", ")
}

func (q Query) String() string {
	res := make([]string, len(q))
	for i, v := range q {
		switch x := v.(type) {
		case string:
			res[i] = x
		case int:
			res[i] = strconv.Itoa(x)
		case Columns:
			res[i] = x.String()
		default:
			panic(fmt.Errorf("the type %T is invalid", x))
		}
	}
	return strings.Join(res, " ")
}

func Set(m interface{}, tag string) {
	v := reflect.ValueOf(m)

	if v.Type().Kind() != reflect.Ptr || v.Type().Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("the model must be pointer to a struct; got %T", v.Type()))
	}

	st := v.Elem()

	if !st.FieldByName("Model").IsValid() {
		panic(fmt.Errorf("the struct %s doesn't embed the type Model", v.Type()))
	}

	columns := st.FieldByName("Columns").Addr().Interface().(*Columns)
	values := st.FieldByName("Values").Addr().Interface().(*Values)

	parseFields(st, tag, columns, values)
}

func parseFields(st reflect.Value, tag string, columns *Columns, values *Values) {
	for i := 0; i < st.NumField(); i++ {
		// check if we are in a embedded struct
		if st.Type().Field(i).Type.Kind() == reflect.Struct {
			parseFields(st.Field(i), tag, columns, values)
		}

		// check if is tagged
		name := st.Type().Field(i).Tag.Get(tag)
		if name == "" {
			continue
		}

		// derive the interface
		inter := st.Field(i).Addr().Interface()

		// change the original input
		*columns = append(*columns, name)
		*values = append(*values, inter)
	}
}
