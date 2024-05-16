package bidb

const (
	opAnd = iota
	opOr
	opAndNot
)

type Result[T any] struct {
	startIdx int

	indexes []int
	ops     []int

	db *DB[T]
}

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

func (res *Result[T]) Get(dest []T) []T {
	if res.db == nil {
		return nil
	}

	res.db.mx.RLock()
	defer res.db.mx.RUnlock()

	r, okSI := res.db.indexes[res.startIdx]
	if !okSI {
		return nil
	}

	result := make([]uint64, len(r))
	copy(result, r)

	if len(res.indexes) == 0 {
		return res.db.indexValues(result, dest)
	}

	var idx [][]uint64
	for _, i := range res.indexes {
		v, ok := res.db.indexes[i]
		if !ok {
			return nil
		}
		idx = append(idx, v)
	}

	for i := 0; i < len(idx); i++ {
		switch res.ops[i] {
		case opAnd:
			var j int

			for {
				if j >= len(result) || j >= len(idx[i]) {
					break
				}

				result[j] &= idx[i][j]
				j++
			}
		case opOr:
			var j int

			for {
				if j >= len(idx[i]) {
					break
				}

				if j >= len(result) {
					result = append(result, idx[i][j:]...)
					break
				}

				result[j] |= idx[i][j]
				j++
			}
		case opAndNot:
			var j int

			for {
				if j >= len(result) || j >= len(idx[i]) {
					break
				}

				result[j] &= ^idx[i][j]
				j++
			}
		}
	}

	return res.db.indexValues(result, dest)
}
