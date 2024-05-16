# BIDB

Bit Indexed DataBase

## Introduction

BIDB is a simple zero-dependencies embedded database that uses bit arrays to store data indexes.
And allows you to query the data using the indexes.

You can only get items by index and by their combinations.
Direct access to the data is not provided.

## Example

```bash
go get github.com/negasus/bidb
```

```go
package main

import (
	"fmt"

	"github.com/negasus/bidb"
)

type User struct {
	ID   int
	Name string
}

const (
	indexMale = iota + 1
	indexFemale
	indexAdult
)

func main() {
	db := bidb.New[User]()

	// You can add batch data
	db.AddBatch([]User{
		{ID: 1, Name: "Andrew"},
		{ID: 1, Name: "John"},
	}, indexMale)

	// or one by one
	db.
		Add(User{ID: 27, Name: "Bot"}).
		Add(User{ID: 2, Name: "Mark"}, indexMale, indexAdult).
		Add(User{ID: 10, Name: "Felix"}, indexMale, indexAdult).
		Add(User{ID: 5, Name: "Mary"}, indexFemale).
		Add(User{ID: 11, Name: "Kate"}, indexFemale, indexAdult).
		Add(User{ID: 10324, Name: "Janny"}, indexFemale)

	resultMaleAdult := db.Index(indexMale).And(indexAdult)
	defer db.ReleaseResult(resultMaleAdult)
	fmt.Printf("Male adult:       %v\n", resultMaleAdult.Get())

	resultFemaleNotAdult := db.Index(indexFemale).AndNot(indexAdult)
	defer db.ReleaseResult(resultFemaleNotAdult)
	fmt.Printf("Female not adult: %v\n", resultFemaleNotAdult.Get())
}
```

## Usage

Create a new database with the type of the data you want to store.

```go
db := bidb.New[User]()
```

Add data to the database. You can add multiple indexes to the data.

> Index with the value of 0 is reserved for all data.

```go
db.Add(User{ID: 1, Name: "John"}, indexMale).
    Add(User{ID: 27, Name: "Bot"}).
    Add(User{ID: 2, Name: "Mark"}, indexMale, indexAdult).
    Add(User{ID: 10, Name: "Felix"}, indexMale, indexAdult).
    Add(User{ID: 5, Name: "Mary"}, indexFemale).
    Add(User{ID: 11, Name: "Kate"}, indexFemale, indexAdult).
    Add(User{ID: 10324, Name: "Janny"}, indexFemale)
```

You can get data with a specific index.

```go
result := db.Index(indexMale)
```

or you can get all data
    
```go
result := db.All()
```

Methods `Index` and `All` returns a `Result` object. You can use the `Get` method to get the data.

```go
data := result.Get()
```

You can also use the `And`, `Or`, `AndNot` methods to combine indexes.

```go
resultMaleAdult := db.Index(indexMale).And(indexAdult)
resultFemaleNotAdult := db.Index(indexFemale).AndNot(indexAdult)
```

You should release the result when you are done with it.

```go
defer db.ReleaseResult(resultMaleAdult)
defer db.ReleaseResult(resultFemaleNotAdult)
```

