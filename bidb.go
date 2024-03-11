package bidb

import (
	"math/bits"
	"sync"
)

const (
	defaultDataLen    = 64
	defaultIndexesLen = 16
)

type index []uint64

// DB is a simple in-memory database
type DB[T any] struct {
	mx   sync.RWMutex
	data []T

	all     index
	indexes map[int]index

	resultPool sync.Pool
}

// New creates a new DB instance
func New[T any]() *DB[T] {
	db := &DB[T]{
		data:    make([]T, 0, defaultDataLen),
		indexes: make(map[int]index, defaultIndexesLen),
	}

	return db
}

// AddBatch adds a batch of items to the database and indexes them
func (db *DB[T]) AddBatch(items []T, indexes ...int) *DB[T] {
	db.mx.Lock()
	defer db.mx.Unlock()

	for _, item := range items {
		db.data = append(db.data, item)
		pos := len(db.data) - 1

		db.all = db.setPos(db.all, pos)

		for _, idx := range indexes {
			db.indexes[idx] = db.setPos(db.indexes[idx], pos)
		}
	}

	return db
}

// Add adds an item to the database and indexes it
func (db *DB[T]) Add(item T, indexes ...int) *DB[T] {
	db.mx.Lock()
	defer db.mx.Unlock()

	db.data = append(db.data, item)
	pos := len(db.data) - 1

	db.all = db.setPos(db.all, pos)

	for _, idx := range indexes {
		db.indexes[idx] = db.setPos(db.indexes[idx], pos)
	}

	return db
}

func (db *DB[T]) setPos(v []uint64, pos int) []uint64 {
	group := pos / 64
	if pos%64 == 0 && pos > 0 {
		group--
	}

	if group >= len(v) {
		v = append(v, make([]uint64, group-len(v)+1)...)
	}

	v[group] |= 1 << (pos - group*64)
	return v
}

func (db *DB[T]) indexValues(values []uint64) []T {
	var res []T

	vv := make([]int, 0, 64)

	for i, v := range values {
		if v == 0 {
			continue
		}

		vv = unpack(v, vv)
		for _, p := range vv {
			elIdx := p + i*64
			if elIdx >= len(db.data) {
				break
			}
			res = append(res, db.data[elIdx])
		}
		vv = vv[:0]
	}

	return res
}

// Index returns a result set for the given index
func (db *DB[T]) Index(index int) *Result[T] {
	db.mx.RLock()
	defer db.mx.RUnlock()

	res := db.acquireResult()
	res.index = append(res.index, db.indexes[index]...)

	return res
}

// All returns a result set for all items in the database
func (db *DB[T]) All() *Result[T] {
	db.mx.RLock()
	defer db.mx.RUnlock()

	res := db.acquireResult()
	res.index = append(res.index, db.all...)

	return res
}

func (db *DB[T]) acquireResult() *Result[T] {
	r, ok := db.resultPool.Get().(*Result[T])
	if !ok {
		r = &Result[T]{
			index: make([]uint64, 0, 16),
			db:    db,
		}
	}

	return r
}

// ReleaseResult releases a result set
func (db *DB[T]) ReleaseResult(res *Result[T]) {
	res.index = res.index[:0]
	db.resultPool.Put(res)
}

func unpack(u uint64, dest []int) []int {
	for u > 0 {
		i := u & -u
		dest = append(dest, bits.TrailingZeros64(i))
		u ^= i
	}

	return dest
}
