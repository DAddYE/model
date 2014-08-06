package model_test

import (
	"fmt"

	"github.com/daddye/model"
)

type User struct {
	Id     int  `sql:"id"`
	Age    int  // not mapped
	Gender byte // not mapped

	// works with embedded structs
	Bio struct {
		Name    string `sql:"name"`
		Surname string `sql:"last_name"`
	}
}

func ExampleModel() {
	user := new(User) // we want a pointer

	// add some data
	user.Id = 9
	user.Age = 30
	user.Gender = 'm'
	user.Bio.Name = "Stan"
	user.Bio.Surname = "Smith"

	// set our model using the our tag
	m := model.New(user, "sql")
	fmt.Println(m.TagNames())

	// check our values
	fmt.Println(m.Values())

	// we can change it
	user.Bio.Name = "Francine"
	fmt.Println(m.Values())

	// we have also a map for an easier access to our struct:
	surname := m.Map()["last_name"]
	fmt.Println(surname)

	// we have access to the reflections
	fmt.Println(m.Reflections())

	// output:
	// [id name last_name]
	// [9 Stan Smith]
	// [9 Francine Smith]
	// Smith
	// [<int Value> Francine Smith]
}

func ExampleQuery() {
	user := new(User) // we want a pointer

	// add some data
	user.Id = 9
	user.Age = 30
	user.Gender = 'm'
	user.Bio.Name = "Stan"
	user.Bio.Surname = "Smith"

	// set our model using the our tag
	m := model.New(user, "sql")

	// we also provide Query, to help dealing with query strings:
	table := "users"
	query := model.Query{"SELECT", m.TagNames(), "FROM", table, "WHERE id=$1"}
	fmt.Println(query)

	// output:
	// SELECT id, name, last_name FROM users WHERE id=$1
}
