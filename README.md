# Model

[![GoDoc](https://godoc.org/github.com/DAddYE/model?status.svg)](https://godoc.org/github.com/DAddYE/model)

General purpose utilities for the go `struct`.

It uses `reflect`, however we cache the result when you allocate the `model`.

## Install

```go
import "github.com/daddye/model"
```

## Usage:

1. Define your `struct`
2. `tag` it. Each tag will map the field to the column name

```go
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
```

create a pointer and add some data:

```go
user := new(User)
user.Id = 9
user.Age = 30
user.Gender = 'm'
user.Bio.Name = "Stan"
user.Bio.Surname = "Smith"
```

now create a model for this struct:

```go
m := model.New(user, "sql")
```

now you can check, tag names:

```go
m.TagNames() // => [id name last_name]
```

values:

```go
m.Values() // => [9 Stan Smith]
```

you're free to change your values anytime:

```go
user.Bio.Name = "Francine"
m.Values() // => [9 Francine Smith]
```

you have access to `reflect.Value`:

```go
m.Reflections()
```

as well interfaces:

```go
m.Interfaces()
```

this, is especially useful when you want to bind `data` to your `struct`, for example, if you're
using `database/sql` you can:

```go
db.QueryRow("SELECT * FROM users WHERE id = $1", 9).Scan(m.Interfaces()...)
```

since query string construction is a common pattern, I have added a type `Query` to deal with it.

So you can rewrite the latter better using our type
[Query](https://godoc.org/github.com/DAddYE/model#Query):

```go
table := "users"
query := model.Query{"SELECT", m.TagNames(), "FROM", table, "WHERE id=$1"}
// query => "SELECT id, name, last_name FROM users WHERE id=$1"

db.QueryRow(query.String(), 9).Scan(m.Interfaces()...)
```

## LICENSE

Copyright (C) 2014 Davide D'Agostino

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
associated documentation files (the "Software"), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge, publish, distribute,
sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial
portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES
OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
