## Model.go

This is a naive, very small (~50 LOC) way to map `something` to a `struct`.

It works well for example with the standard package `database/sql`

### Usage:

Define a type, `tag` it and embed the `Model` type.

```go
import "github.com/daddye/model"

type Feed struct {
	Id         int    `sql:"id"`
	SourceType string `sql:"source_type"`
	Model
}
```

Finally setup your new model:

```go
f := new(Feed)
model.Set(f, "sql") // the tag you chose before
```

Now you can use it like:

```go
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
```

### LICENSE

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
