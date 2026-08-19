package itemset

import "testing"

func TestUnion(t *testing.T) {
	a := New("a", "b", "c")
	b := New("b", "d")
	got := Union(a, b)
	want := Set{"a", "b", "c", "d"}
	if !Equal(got, want) {
		t.Errorf("Union = %v, want %v", got, want)
	}
}

func TestIntersect(t *testing.T) {
	a := New("a", "b", "c")
	b := New("b", "c", "d")
	got := Intersect(a, b)
	want := Set{"b", "c"}
	if !Equal(got, want) {
		t.Errorf("Intersect = %v, want %v", got, want)
	}
}

func TestDifference(t *testing.T) {
	a := New("a", "b", "c")
	b := New("b")
	got := Difference(a, b)
	want := Set{"a", "c"}
	if !Equal(got, want) {
		t.Errorf("Difference = %v, want %v", got, want)
	}
}

func TestIsSubset(t *testing.T) {
	a := New("b", "c")
	b := New("a", "b", "c", "d")
	if !IsSubset(a, b) {
		t.Error("expected a to be subset of b")
	}
	if IsSubset(b, a) {
		t.Error("b should not be subset of a")
	}
}

func TestSubsets(t *testing.T) {
	s := New("a", "b", "c")
	subs := Subsets(s, 2)
	if len(subs) != 3 {
		t.Errorf("C(3,2) = %d, want 3", len(subs))
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	s := New("apple", "banana", "cherry")
	data := Encode(s)
	got, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if !Equal(got, s) {
		t.Errorf("roundtrip: got %v, want %v", got, s)
	}
}
