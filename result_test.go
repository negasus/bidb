package bidb

import "testing"

type Item struct {
	ID int
}

func TestResultAnd(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)
	db.Add(Item{ID: 2}, 1, 2)
	db.Add(Item{ID: 3}, 1, 2)
	db.Add(Item{ID: 4}, 2)
	db.Add(Item{ID: 5}, 3)

	res := db.Index(1).And(2).Get()

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 2 {
		t.Errorf("Expected 2, got %d", res[0].ID)
	}

	if res[1].ID != 3 {
		t.Errorf("Expected 3, got %d", res[1].ID)
	}
}

func TestResultOr(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)
	db.Add(Item{ID: 2}, 1, 2)
	db.Add(Item{ID: 3}, 1, 2, 4)
	db.Add(Item{ID: 4}, 2)
	db.Add(Item{ID: 5}, 3)

	res := db.Index(4).Or(3).Get()

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 3 {
		t.Errorf("Expected 3, got %d", res[0].ID)
	}

	if res[1].ID != 5 {
		t.Errorf("Expected 5, got %d", res[1].ID)
	}
}

func TestResultOr_secondIndexBig(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)

	for i := 2; i < 100; i++ {
		db.Add(Item{ID: i}, 2)
	}
	db.Add(Item{ID: 1000}, 3)

	res := db.Index(1).Or(3).Get()

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 1 {
		t.Errorf("Expected 1, got %d", res[0].ID)
	}

	if res[1].ID != 1000 {
		t.Errorf("Expected 1000, got %d", res[1].ID)
	}
}

func TestAndNot_ManyItems(t *testing.T) {
	db := New[Item]()

	for i := 0; i < 100; i++ {
		db.Add(Item{ID: i}, 1)
	}

	db.Add(Item{ID: 1000}, 2)

	res := db.All().AndNot(1).Get()

	if len(res) != 1 {
		t.Errorf("Expected 1, got %d", len(res))
	}

	if res[0].ID != 1000 {
		t.Errorf("Expected 1000, got %d", res[0].ID)
	}
}

func TestResultOr_firstIndexBig(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)

	for i := 2; i < 100; i++ {
		db.Add(Item{ID: i}, 2)
	}
	db.Add(Item{ID: 1000}, 3)

	res := db.Index(3).Or(1).Get()

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 1 {
		t.Errorf("Expected 1, got %d", res[0].ID)
	}

	if res[1].ID != 1000 {
		t.Errorf("Expected 1000, got %d", res[1].ID)
	}
}

func TestResultAndNot(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)
	db.Add(Item{ID: 2}, 1, 2)
	db.Add(Item{ID: 3}, 1, 2, 4)
	db.Add(Item{ID: 4}, 2)
	db.Add(Item{ID: 5}, 3)

	res := db.Index(1).AndNot(4).Get()

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 1 {
		t.Errorf("Expected 1, got %d", res[0].ID)
	}

	if res[1].ID != 2 {
		t.Errorf("Expected 2, got %d", res[1].ID)
	}
}

func TestResultNot(t *testing.T) {
	db := New[Item]()
	db.Add(Item{ID: 1}, 1)
	db.Add(Item{ID: 2}, 1, 2)
	db.Add(Item{ID: 3}, 1, 2, 4)
	db.Add(Item{ID: 4}, 2)
	db.Add(Item{ID: 5}, 3)

	res := db.Index(1).Not().Get()

	if len(res) != 2 {
		t.Errorf("Expected 2, got %d", len(res))
	}

	if res[0].ID != 4 {
		t.Errorf("Expected 4, got %d", res[0].ID)
	}

	if res[1].ID != 5 {
		t.Errorf("Expected 5, got %d", res[1].ID)
	}
}
