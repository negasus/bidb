package bidb

// Result is a result of a query
type Result[T any] struct {
	startIdx int

	indexes []int
	ops     []int

	db *DB[T]
}

const (
	opAnd = iota
	opOr
	opAndNot
	opNot
)

func (res *Result[T]) And(index int) *Result[T] {
	res.indexes = append(res.indexes, index)
	res.ops = append(res.ops, opAnd)
	return res
}

func (res *Result[T]) Or(index int) *Result[T] {
	res.indexes = append(res.indexes, index)
	res.ops = append(res.ops, opOr)
	return res
}

func (res *Result[T]) AndNot(index int) *Result[T] {
	res.indexes = append(res.indexes, index)
	res.ops = append(res.ops, opAndNot)
	return res
}

func (res *Result[T]) Not(index int) *Result[T] {
	res.indexes = append(res.indexes, index)
	res.ops = append(res.ops, opNot)
	return res
}

// Get returns the result set
func (res *Result[T]) Get() []T {
	res.db.mx.RLock()

	si, ok := res.db.indexes[res.startIdx]
	if !ok {
		return nil
	}

	if len(res.indexes) == 0 {
		res.db.mx.RUnlock()
		return res.db.indexValues(si)
	}

	var idx [][]uint64
	for _, i := range res.indexes {
		v, oki := res.db.indexes[i]
		if !oki {
			return nil
		}
		idx = append(idx, v)
	}
	res.db.mx.RUnlock()

	for i := 0; i < len(idx); i++ {
		switch res.ops[i] {
		case opAnd:
			for {
				if i >= len(res.index) || i >= len(ai) {
					break
				}

				res.index[i] &= ai[i]
				i++
			}
		case opOr:
			res.or(i)
		case opAndNot:
			res.andNot(i)
		case opNot:
			res.not()
		}
	}

	return nil
}

func (res *Result[T]) and(index int) *Result[T] {
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

func (res *Result[T]) or(index int) *Result[T] {
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

func (res *Result[T]) andNot(index int) *Result[T] {
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

func (res *Result[T]) not() *Result[T] {
	for i := range res.index {
		res.index[i] = ^res.index[i]
	}
	return res
}
