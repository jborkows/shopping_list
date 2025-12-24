package products

import "testing"

func TestNormalizeGroupName(t *testing.T) {
	if _, err := NormalizeGroupName("   "); err == nil {
		t.Fatalf("expected error")
	}
	got, err := NormalizeGroupName(" vegetables ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "vegetables" {
		t.Fatalf("expected trimmed name, got %q", got)
	}
}

func TestNormalizeProductName(t *testing.T) {
	if _, err := NormalizeProductName("   "); err == nil {
		t.Fatalf("expected error")
	}
	got, err := NormalizeProductName(" carrots ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "carrots" {
		t.Fatalf("expected trimmed name, got %q", got)
	}
}
