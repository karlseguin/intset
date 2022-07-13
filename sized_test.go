package intset

import (
	"math/rand"
	"testing"
)

func Test_Sized_SetsAValue(t *testing.T) {
	s := NewSized(20)
	for i := 0; i < 30; i++ {
		s.Set(i)
		AssertTrue(t, s.Exists(i))
	}
	for i := 0; i < 30; i++ {
		AssertTrue(t, s.Exists(i))
	}
}

func Test_Sized_Exists(t *testing.T) {
	s := NewSizedConfig(20, NewConfig())
	for i := 0; i < 100; i++ {
		AssertFalse(t, s.Exists(i))
		s.Set(i)
	}
	for i := 0; i < 100; i++ {
		AssertTrue(t, s.Exists(i))
	}
}

func Test_Sized_SizeLessThanBucket(t *testing.T) {
	config := NewConfig().BucketSize(8).BucketGrowBy(4)
	s := NewSizedConfig(7, config)
	s.Set(32)
	AssertTrue(t, s.Exists(32))
	AssertFalse(t, s.Exists(33))
}

func Test_Sized_RemoveNonMembers(t *testing.T) {
	s := NewSized(100)
	AssertFalse(t, s.Remove(329))
}

func Test_Sized_RemovesMembers(t *testing.T) {
	s := NewSized(100)
	for i := 0; i < 10; i++ {
		s.Set(i)
	}
	AssertFalse(t, s.Remove(20))
	AssertTrue(t, s.Remove(2))
	AssertFalse(t, s.Remove(2))
	AssertFalse(t, s.Exists(2))
	AssertEqual(t, s.Len(), 9)
}

func Test_Sized_IntersectsTwoSets(t *testing.T) {
	s1 := NewSized(10)
	s2 := NewSized(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Intersect([]Set{s1, s2})
	AssertFalse(t, s.Exists(1))
	AssertTrue(t, s.Exists(2))
	AssertTrue(t, s.Exists(3))
	AssertFalse(t, s.Exists(4))
	AssertFalse(t, s.Exists(5))
}

func Test_Sized_UnionsTwoSets(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s1 := NewSized(10)
		s2 := NewSized(10)
		s1.Set(1)
		s1.Set(2)
		s1.Set(3)

		s2.Set(2)
		s2.Set(3)
		s2.Set(4)

		s := Union([]Set{s1, s2})
		AssertTrue(t, s.Exists(1))
		AssertTrue(t, s.Exists(2))
		AssertTrue(t, s.Exists(3))
		AssertTrue(t, s.Exists(4))
		AssertFalse(t, s.Exists(5))
	}
}

func Benchmark_SizedPopulate(b *testing.B) {
	s := NewSized(10000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(i % 10000000)
	}
}

func Benchmark_SizedDenseExists(b *testing.B) {
	s := NewSized(1000000)
	for i := 0; i < 1000000; i++ {
		s.Set(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exists(i % 1000000)
	}
}

func Benchmark_SizedSparseExists(b *testing.B) {
	s := NewSized(1000000)
	for i := 0; i < 1000000; i++ {
		if i%10 == 0 {
			s.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exists(i % 1000000)
	}
}

func Benchmark_SizedDenseIntersect(b *testing.B) {
	s1 := NewSized(100000)
	for i := 0; i < 100000; i++ {
		if rand.Intn(10) != 0 {
			s1.Set(i)
		}
	}
	s2 := NewSized(1000)
	for i := 0; i < 1000; i++ {
		if rand.Intn(10) != 0 {
			s2.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Intersect([]Set{s1, s2})
	}
}

// Benchmarks for map[int]struct{}

func Benchmark_SizedMapDenseExists(b *testing.B) {
	s := make(map[int]struct{})
	for i := 0; i < 1000000; i++ {
		s[i] = struct{}{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s[i]
	}
}

func Benchmark_SizedMapSparseExists(b *testing.B) {
	s := make(map[int]struct{})
	for i := 0; i < 1000000; i++ {
		if i%10 == 0 {
			s[i] = struct{}{}
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s[i%1000000]
	}
}

// Two values are equal
func AssertEqual[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("\nexpected: '%v'\nto equal: '%v'", actual, expected)
		t.FailNow()
	}
}

// A value is true
func AssertTrue(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Error("expected true, got false")
		t.FailNow()
	}
}

// A value is false
func AssertFalse(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Error("expected false, got true")
		t.FailNow()
	}
}
