package bidb

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	db := New[Item]()

	db.Add(Item{ID: 1}, 1)

	if len(db.data) != 1 {
		t.Errorf("Expected 1, got %d", len(db.data))
	}

	if db.data[0].ID != 1 {
		t.Errorf("Expected 1, got %d", db.data[0].ID)
	}

	if len(db.indexes) != 2 {
		t.Errorf("Expected 2, got %d", len(db.indexes))
	}

	if len(db.indexes[1]) != 1 {
		t.Errorf("Expected 1, got %d", len(db.indexes[1]))
	}

	if db.indexes[1][0] != 1 {
		t.Errorf("Expected 1, got %d", db.indexes[1][0])
	}

	if len(db.indexes[0]) != 1 {
		t.Errorf("Expected 1, got %d", len(db.indexes[0]))
	}

	if db.indexes[0][0] != 1 {
		t.Errorf("Expected 1, got %d", db.indexes[0][0])
	}

}

func TestAddBatch(t *testing.T) {
	db := New[Item]()

	db.AddBatch([]Item{{ID: 1}, {ID: 2}}, 1)

	if len(db.data) != 2 {
		t.Errorf("Expected 2, got %d", len(db.data))
	}

	if db.data[0].ID != 1 {
		t.Errorf("Expected 1, got %d", db.data[0].ID)
	}
	if db.data[1].ID != 2 {
		t.Errorf("Expected 2, got %d", db.data[1].ID)
	}

	if len(db.indexes) != 2 {
		t.Errorf("Expected 2, got %d", len(db.indexes))
	}

	if len(db.indexes[1]) != 1 {
		t.Errorf("Expected 1, got %d", len(db.indexes[1]))
	}

	if db.indexes[1][0] != 3 {
		t.Errorf("Expected 3, got %d", db.indexes[1][0])
	}

	if len(db.indexes[0]) != 1 {
		t.Errorf("Expected 1, got %d", len(db.indexes[0]))
	}

	if db.indexes[0][0] != 3 {
		t.Errorf("Expected 3, got %d", db.indexes[0][0])
	}
}

func Test_setPos(t *testing.T) {
	db := New[Item]()

	var v []uint64

	v = db.setPos([]uint64{0}, 0)
	checkPos(t, v, []uint64{0x00_00_00_00_00_00_00_01})

	v = db.setPos([]uint64{0x00_00_00_00_00_00_00_01}, 5)
	checkPos(t, v, []uint64{0x00_00_00_00_00_00_00_21})

	v = db.setPos([]uint64{0}, 5)
	checkPos(t, v, []uint64{0x00_00_00_00_00_00_00_20})

	v = db.setPos([]uint64{0}, 63)
	checkPos(t, v, []uint64{0x80_00_00_00_00_00_00_00})

	v = db.setPos([]uint64{0}, 64)
	checkPos(t, v, []uint64{0x00_00_00_00_00_00_00_00, 0x00_00_00_00_00_00_00_01})

	v = db.setPos([]uint64{0}, 128)
	checkPos(t, v, []uint64{0x00_00_00_00_00_00_00_00, 0x00_00_00_00_00_00_00_00, 0x00_00_00_00_00_00_00_01})
}

func checkPos(t *testing.T, v, expected []uint64) {
	if !reflect.DeepEqual(v, expected) {
		t.Fatalf("Expected %.16x, got %.16x", expected, v)
	}
}

func Test_indexValues(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)
	db.Add(Item{ID: 2}, 2)
	db.Add(Item{ID: 3}, 1)

	res := db.indexValues([]uint64{5}, nil) // 0000 0101

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 1 {
		t.Errorf("Expected 1, got %d", res[0].ID)
	}

	if res[1].ID != 3 {
		t.Errorf("Expected 3, got %d", res[1].ID)
	}
}

func Test_indexValues_limit(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)

	res := db.indexValues([]uint64{7}, nil) // 0000 0111

	if len(res) != 1 {
		t.Errorf("Expected 1, got %d", len(res))
	}

	if res[0].ID != 1 {
		t.Errorf("Expected 1, got %d", res[0].ID)
	}
}

func TestIndex(t *testing.T) {
	db := New[Item]()

	db.Add(Item{ID: 1}, 1)
	db.Add(Item{ID: 2}, 2)
	db.Add(Item{ID: 3}, 1)

	res := db.Index(1).Get(nil)

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 1 {
		t.Errorf("Expected 1, got %d", res[0].ID)
	}

	if res[1].ID != 3 {
		t.Errorf("Expected 3, got %d", res[1].ID)
	}
}

func Test_unpack(t *testing.T) {
	type args struct {
		u    uint64
		dest []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "empty",
			args: args{
				u:    0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				dest: []int{},
			},
			want: []int{},
		},
		{
			name: "0-3-10-18-33",
			args: args{
				u:    0b00000000_00000000_00000000_00000010_00000000_00000100_00000100_00001001,
				dest: []int{},
			},
			want: []int{0, 3, 10, 18, 33},
		},
		{
			name: "0-63",
			args: args{
				u:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				dest: []int{},
			},
			want: []int{0, 63},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unpack(tt.args.u, tt.args.dest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpack() = %v, want %v", got, tt.want)
			}
		})
	}
}
