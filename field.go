package model

import "reflect"

// Field represent the field of a struct, is an extension of `reflect.StructField` and adds
// few convenient methods.
type Field struct {
	TagName   string
	Interface interface{}
	reflect.StructField
}

func fields(sv reflect.Value, tag string) ([]reflect.Value, []reflect.StructField) {
	v := make([]reflect.Value, 0)
	t := make([]reflect.StructField, 0)
	st := sv.Type()

	for i := 0; i < st.NumField(); i++ {
		if st.Field(i).Tag.Get(tag) == "-" {
			continue
		}

		// walk inside an embedded struct if it has no tag.
		if st.Field(i).Type.Kind() == reflect.Struct && st.Field(i).Tag.Get(tag) == "" {
			vn, tn := fields(sv.Field(i), tag)
			v = append(v, vn...)
			t = append(t, tn...)
			continue
		}

		if st.Field(i).Tag.Get(tag) != "" {
			v = append(v, sv.Field(i))
			t = append(t, st.Field(i))
		}
	}

	if len(v) != len(t) {
		panic("internal error")
	}

	return v, t
}
