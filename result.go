package bidb

// Result is a result of a query
type Result[T any] struct {
	index []uint64
	db    *DB[T]
}

// And combines the current result with another index
func (res *Result[T]) And(index int) *Result[T] {
	var i int

	ai := res.db.indexes[index]

	for {
		if i >= len(res.index) || i >= len(ai) {
			break
		}

		res.index[i] &= ai[i]
		i++
	}
	return res
}

// Or combines the current result with another index
func (res *Result[T]) Or(index int) *Result[T] {
	var i int

	ai := res.db.indexes[index]

	for {
		if i >= len(ai) {
			break
		}

		if i >= len(res.index) {
			res.index = append(res.index, ai[i:]...)
			break
		}

		res.index[i] |= ai[i]
		i++
	}
	return res
}

// AndNot combines the current result with another index
func (res *Result[T]) AndNot(index int) *Result[T] {
	var i int

	ai := res.db.indexes[index]

	for {
		if i >= len(res.index) || i >= len(ai) {
			break
		}

		res.index[i] &= ^ai[i]
		i++
	}
	return res
}

// Not inverts the current result
func (res *Result[T]) Not() *Result[T] {
	for i := range res.index {
		res.index[i] = ^res.index[i]
	}
	return res
}

// Get returns the result set
func (res *Result[T]) Get() []T {
	return res.db.indexValues(res.index)
}
