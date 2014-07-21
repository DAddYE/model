package model

import (
	"fmt"
	"reflect"
)

/*
  Usage:

  Define a type, tag it and embed the Model type.

	type Feed struct {
		Id         int    `sql:"id"`
		SourceType string `sql:"source_type"`
		Model
	}

  Finally setup your new model:

	f := new(Feed)
	SetModel(f, "sql")

  Now you can use it like:

	db, err := sql.Open("postgres", "...")
	rows, err := db.Query("SELECT " + strings.Join(f.Columns, ", ") + " FROM feeds LIMIT 10")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err := rows.Scan(f.Values...)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#+v\n", f.SourceType)
	}
*/
type Model struct {
	Err     error
	Columns []string
	Values  []interface{}
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

	columns := st.FieldByName("Columns").Addr().Interface().(*[]string)
	values := st.FieldByName("Values").Addr().Interface().(*[]interface{})

	for i := 0; i < st.NumField(); i++ {
		name := st.Type().Field(i).Tag.Get(tag)
		if name == "" {
			continue
		}
		inter := st.Field(i).Addr().Interface()
		*columns = append(*columns, name)
		*values = append(*values, inter)
	}
}
