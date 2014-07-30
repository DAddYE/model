# Model

This is a naive, very small way to map `something` to a `struct`.

It works well for example with the standard package `database/sql`, but also others like
`cassandra/cql` should be fine.

It uses `reflections`, but only **once**, in particular during the allocation time.

It is very cheap because as the standard library, you allocate **once** an then you'll **reuse**
your `var` to bind the result.

Instead to use the `interface` approach in order to avoid type assertions or reflection (when
binding the result) this package use `embedded structs` to achieve a similar result.

## Install

```go
import "github.com/daddye/model"
```

## Usage:

### Definition

1. Define your model type (Usually the name of a Table)
2. `tag` it. Each tag will map the field to the column name
3. Embed the `Model` type.

```go
import "github.com/daddye/model"

type Feed struct {
	Id               int    `sql:"id"`
	SoruceFileFormat string `sql:"source_file_format"`
	SourceType       string `sql:"source_type"`
	Ftp              Ftp
	model.Model      // <<<< this
}
```

You can embed many structs as you want in this case `FTP` is a separate one:

```go
type Ftp struct {
	FileName string `sql:"source_ftp_file_name"`
	Username string `sql:"source_ftp_username"`
}
```

### Allocation

You must provide an allocation function in order to allow `Model` to derive/map the struct's fields
to the sql counterpart:

```go
func NewFeed(db *sql.DB) (f *Feed) {
	f = new(Feed)
	f.Interface = db
	f.Table = "feeds"
	model.Set(f, "sql")
	return
}
```

Here you have done few things:

1. Allocation of the struct (as usual)
2. Map the struct with a `connection`
3. Set the table where the struct `Feed` will refer to.
4. Map the struct's fields with their tags and setup a `Value` to use later.

### Interface

At this point you need to define the `interface` (it's not really that and it's not mandatory) so you
may want to fetch `one` or `multi` row(s), you can now define your functions as:

```go
func (f *Feed) First(conditions string, args ...interface{}) error {
	m := &f.Model
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions, "LIMIT 1"}

	m.Err = m.Interface.(*sql.DB).QueryRow(query.String(), args...).Scan(m.Values...)
	return m.Err
}

func (f *Feed) Find(conditions string, args ...interface{}) (*model.Iter, error) {
	m := &f.Model
	query := model.Query{"SELECT", m.Columns, "FROM", m.Table, conditions}

	rows, err := m.Interface.(*sql.DB).Query(query.String(), args...)
	if err != nil {
		m.Err = err
		return nil, err
	}

	iter := model.Iter{Model: m, Rows: rows}
	return &iter, nil
}
```

_SIDE NOTE: What is `Query`? Nothing complicated, just a little helper to `Join(something, ", ")`,
check it out on [model.go](/model.go)_

Since it's a quite common pattern we provide (and you can do it too) `generics` functions. So for
`database/sql` you can import [query](/model/sql) and write just that:

```go
import query "github.com/daddye/model/sql"

func (f *Feed) First(conditions string, args ...interface{}) error {
	return query.First(&f.Model, conditions, args...)
}

func (f *Feed) Find(conditions string, args ...interface{}) (*model.Iter, error) {
	return query.Find(&f.Model, conditions, args...)
}
```

### Usage

Now you can use it on a single row:

```go
db, err := sql.Open("postgres", "user=... dbname=... sslmode=disable")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

feed := NewFeed(db)
err = feed.First("WHERE ID = $1", 123)
if err != nil {
	fmt.Print(err)
	return
}
fmt.Printf("Id: %d, FileName: %s\n", feed.Id, feed.Ftp.Username)
```

For multiple results instead you can do:

```go
db, err := sql.Open("postgres", "user=... dbname=... sslmode=disable")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

feed := NewFeed(db)
rows, err := feed.Find("WHERE source_type = 'ftp'")
if err != nil {
	fmt.Print(err)
	return
}
for rows.Next() {
	fmt.Printf("Id: %d, FileName: %s\n", feed.Id, feed.Ftp.Username)
}
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
