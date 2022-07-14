package intset

import "testing"

// Common testing utility functions

// AssertEqual checks if two values are equal
func AssertEqual[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("\nexpected: '%v'\nto equal: '%v'", actual, expected)
		t.FailNow()
	}
}

// AssertTrue checks if a value is true
func AssertTrue(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Error("expected true, got false")
		t.FailNow()
	}
}

// AssertFalse checks if a value is false
func AssertFalse(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Error("expected false, got true")
		t.FailNow()
	}
}
