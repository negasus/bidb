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
	fmt.Printf("Male adult:       %v\n", resultMaleAdult.Get(nil))

	resultFemaleNotAdult := db.Index(indexFemale).AndNot(indexAdult)
	defer db.ReleaseResult(resultFemaleNotAdult)
	fmt.Printf("Female not adult: %v\n", resultFemaleNotAdult.Get(nil))

	// resultNotAdult := db.Index(indexAdult).Not()
	// defer db.ReleaseResult(resultNotAdult)
	// fmt.Printf("Not adult:        %v\n", resultNotAdult.Get())
}
