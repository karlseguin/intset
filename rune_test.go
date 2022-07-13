package intset

import (
	"math/rand"
	"testing"
)

func Test_Rune_SetsAValue(t *testing.T) {
	s := NewRune(20)
	for i := rune(0); i < 30; i++ {
		s.Set(i)
		AssertTrue(t, s.Exists(i))
	}
	for i := rune(0); i < 30; i++ {
		AssertTrue(t, s.Exists(i))
	}
}

func Test_Rune_Exists(t *testing.T) {
	s := NewRune(20)
	for i := rune(0); i < 10; i++ {
		AssertFalse(t, s.Exists(i))
		s.Set(i)
	}
	for i := rune(0); i < 10; i++ {
		AssertTrue(t, s.Exists(i))
	}
}

func Test_Rune_SizeLessThanBucket(t *testing.T) {
	s := NewRune(rune(Default.bucketSize) - 1)
	s.Set(32)
	AssertTrue(t, s.Exists(32))
	AssertFalse(t, s.Exists(33))
}

func Test_Rune_RemoveNonMembers(t *testing.T) {
	s := NewRune(100)
	AssertFalse(t, s.Remove(329))
}

func Test_Rune_RemovesMembers(t *testing.T) {
	s := NewRune(100)
	for i := rune(0); i < 10; i++ {
		s.Set(i)
	}
	AssertFalse(t, s.Remove(20))
	AssertTrue(t, s.Remove(2))
	AssertFalse(t, s.Remove(2))
	AssertFalse(t, s.Exists(2))
	AssertEqual(t, s.Len(), 9)
}

func Test_Rune_IntersectsTwoSets(t *testing.T) {
	s1 := NewRune(10)
	s2 := NewRune(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := IntersectRune([]SetRune{s1, s2})
	AssertFalse(t, s.Exists(1))
	AssertTrue(t, s.Exists(2))
	AssertTrue(t, s.Exists(3))
	AssertFalse(t, s.Exists(4))
	AssertFalse(t, s.Exists(5))
}

func Test_Rune_UnionsTwoSets(t *testing.T) {
	s1 := NewRune(10)
	s2 := NewRune(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := UnionRune([]SetRune{s1, s2})
	AssertTrue(t, s.Exists(1))
	AssertTrue(t, s.Exists(2))
	AssertTrue(t, s.Exists(3))
	AssertTrue(t, s.Exists(4))
	AssertFalse(t, s.Exists(5))
}

func Benchmark_RunePopulate(b *testing.B) {
	s := NewRune(10000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(rune(i % 10000000))
	}
}

func Benchmark_RuneDenseExists(b *testing.B) {
	s := NewRune(1000000)
	for i := rune(0); i < 1000000; i++ {
		s.Set(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exists(rune(i % 1000000))
	}
}

func Benchmark_RuneSparseExists(b *testing.B) {
	s := NewRune(1000000)
	for i := rune(0); i < 1000000; i++ {
		if i%10 == 0 {
			s.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exists(rune(i % 1000000))
	}
}

func Benchmark_RuneDenseIntersect(b *testing.B) {
	s1 := NewRune(100000)
	for i := rune(0); i < 100000; i++ {
		if rand.Intn(10) != 0 {
			s1.Set(i)
		}
	}
	s2 := NewRune(1000)
	for i := rune(0); i < 1000; i++ {
		if rand.Intn(10) != 0 {
			s2.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntersectRune([]SetRune{s1, s2})
	}
}

// Benchmarks for map[rune]struct{}

func Benchmark_RuneMapDenseExists(b *testing.B) {
	s := make(map[rune]struct{})
	for i := rune(0); i < 1000000; i++ {
		s[i] = struct{}{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s[rune(i)]
	}
}

func Benchmark_RuneMapSparseExists(b *testing.B) {
	s := make(map[rune]struct{})
	for i := rune(0); i < 1000000; i++ {
		if i%10 == 0 {
			s[i] = struct{}{}
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s[rune(i%1000000)]
	}
}
