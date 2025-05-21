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
	mx         sync.RWMutex
	data       []T
	indexes    map[int]index
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

// Reset removes all items from the database and resets indexes
func (db *DB[T]) Reset() {
	db.mx.Lock()
	defer db.mx.Unlock()

	db.data = db.data[:0]
	for k := range db.indexes {
		delete(db.indexes, k)
	}
}

// AddBatch adds a batch of items to the database and indexes them
func (db *DB[T]) AddBatch(items []T, indexes ...int) *DB[T] {
	db.mx.Lock()
	defer db.mx.Unlock()

	for _, item := range items {
		db.data = append(db.data, item)
		pos := len(db.data) - 1

		db.indexes[0] = db.setPos(db.indexes[0], pos)

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

	db.indexes[0] = db.setPos(db.indexes[0], pos)

	for _, idx := range indexes {
		db.indexes[idx] = db.setPos(db.indexes[idx], pos)
	}

	return db
}

// FillFrom fills the current database with data from another database
func (db *DB[T]) FillFrom(src *DB[T]) {
	// reset
	db.mx.Lock()
	defer db.mx.Unlock()

	db.data = db.data[:0]
	for k := range db.indexes {
		delete(db.indexes, k)
	}

	// copy
	for i := 0; i < len(src.data); i++ {
		db.data = append(db.data, src.data[i])
	}
	for k, v := range src.indexes {
		db.indexes[k] = v
	}
}

func (db *DB[T]) setPos(v []uint64, pos int) []uint64 {
	group := pos / 64

	if group >= len(v) {
		v = append(v, make([]uint64, group-len(v)+1)...)
	}

	shift := pos - group*64

	v[group] |= 1 << shift

	return v
}

func (db *DB[T]) indexValues(values []uint64, dest []T) []T {
	db.mx.RLock()
	defer db.mx.RUnlock()

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
			dest = append(dest, db.data[elIdx])
		}
		vv = vv[:0]
	}

	return dest
}

func (db *DB[T]) Index(index int) *Result[T] {
	res := db.acquireResult()
	res.startIdx = index

	return res
}

func (db *DB[T]) All() *Result[T] {
	res := db.acquireResult()
	res.startIdx = 0

	return res
}

func (db *DB[T]) acquireResult() *Result[T] {
	r, ok := db.resultPool.Get().(*Result[T])
	if !ok {
		r = &Result[T]{
			db: db,
		}
	}

	return r
}

// ReleaseResult releases a result set
func (db *DB[T]) ReleaseResult(res *Result[T]) {
	res.startIdx = 0
	res.ops = res.ops[:0]
	res.indexes = res.indexes[:0]
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
