package shoppinglist

import "testing"

func TestNormalizeItemName(t *testing.T) {
	got, err := NormalizeItemName("  milk   2%  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "milk 2%" {
		t.Fatalf("got %q, want %q", got, "milk 2%")
	}
}

func TestNormalizeItemName_Empty(t *testing.T) {
	_, err := NormalizeItemName("   ")
	if err == nil {
		t.Fatalf("expected error")
	}
}
