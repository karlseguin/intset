package intset

import (
	"math/rand"
	"testing"

	expect "github.com/karlseguin/expect"
)

type RuneTest struct{}

func Test_Rune(t *testing.T) {
	expect.Expectify(new(RuneTest), t)
}

func (RuneTest) SetsAValue() {
	s := NewRune(20)
	for i := rune(0); i < 30; i++ {
		s.Set(i)
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
	for i := rune(0); i < 30; i++ {
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
}

func (RuneTest) Exists() {
	s := NewRune(20)
	for i := rune(0); i < 10; i++ {
		expect.Expect(s.Exists(i)).To.Equal(false)
		s.Set(i)
	}
	for i := rune(0); i < 10; i++ {
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
}

func (RuneTest) SizeLessThanBucket() {
	s := NewRune(rune(bucketSize) - 1)
	s.Set(32)
	expect.Expect(s.Exists(32)).To.Equal(true)
	expect.Expect(s.Exists(33)).To.Equal(false)
}

func (RuneTest) RemoveNonMembers() {
	s := NewRune(100)
	expect.Expect(s.Remove(329)).To.Equal(false)
}

func (RuneTest) RemovesMembers() {
	s := NewRune(100)
	for i := rune(0); i < 10; i++ {
		s.Set(i)
	}
	expect.Expect(s.Remove(20)).To.Equal(false)
	expect.Expect(s.Remove(2)).To.Equal(true)
	expect.Expect(s.Remove(2)).To.Equal(false)
	expect.Expect(s.Exists(2)).To.Equal(false)
	expect.Expect(s.Len()).To.Equal(9)
}

func (RuneTest) IntersectsTwoSets() {
	s1 := NewRune(10)
	s2 := NewRune(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := IntersectRune([]Rune{s1, s2})
	expect.Expect(s.Exists(1)).To.Equal(false)
	expect.Expect(s.Exists(2)).To.Equal(true)
	expect.Expect(s.Exists(3)).To.Equal(true)
	expect.Expect(s.Exists(4)).To.Equal(false)
	expect.Expect(s.Exists(5)).To.Equal(false)
}

func (RuneTest) UnionsTwoSets() {
	s1 := NewRune(10)
	s2 := NewRune(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := UnionRune([]Rune{s1, s2})
	expect.Expect(s.Exists(1)).To.Equal(true)
	expect.Expect(s.Exists(2)).To.Equal(true)
	expect.Expect(s.Exists(3)).To.Equal(true)
	expect.Expect(s.Exists(4)).To.Equal(true)
	expect.Expect(s.Exists(5)).To.Equal(false)
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
		IntersectRune([]Rune{s1, s2})
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
