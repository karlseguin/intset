package intset

import (
	"math/rand"
	"testing"

	expect "github.com/karlseguin/expect"
)

type SizedTest struct{}

func Test_Sized(t *testing.T) {
	expect.Expectify(new(SizedTest), t)
}

func (SizedTest) SetsAValue() {
	s := NewSized(20)
	for i := 0; i < 30; i++ {
		s.Set(i)
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
	for i := 0; i < 30; i++ {
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
}

func (SizedTest) Exists() {
	s := NewSized(20)
	for i := 0; i < 100; i++ {
		expect.Expect(s.Exists(i)).To.Equal(false)
		s.Set(i)
	}
	for i := 0; i < 100; i++ {
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
}

func (SizedTest) SizeLessThanBucket() {
	s := NewSized(bucketSize - 1)
	s.Set(32)
	expect.Expect(s.Exists(32)).To.Equal(true)
	expect.Expect(s.Exists(33)).To.Equal(false)
}

func (SizedTest) RemoveNonMembers() {
	s := NewSized(100)
	expect.Expect(s.Remove(329)).To.Equal(false)
}

func (SizedTest) RemovesMembers() {
	s := NewSized(100)
	for i := 0; i < 10; i++ {
		s.Set(i)
	}
	expect.Expect(s.Remove(20)).To.Equal(false)
	expect.Expect(s.Remove(2)).To.Equal(true)
	expect.Expect(s.Remove(2)).To.Equal(false)
	expect.Expect(s.Exists(2)).To.Equal(false)
	expect.Expect(s.Len()).To.Equal(9)
}

func (SizedTest) IntersectsTwoSets() {
	s1 := NewSized(10)
	s2 := NewSized(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Intersect([]intSet{s1, s2})
	expect.Expect(s.Exists(1)).To.Equal(false)
	expect.Expect(s.Exists(2)).To.Equal(true)
	expect.Expect(s.Exists(3)).To.Equal(true)
	expect.Expect(s.Exists(4)).To.Equal(false)
	expect.Expect(s.Exists(5)).To.Equal(false)
}

func (SizedTest) UnionsTwoSets() {
	for i := 0; i < 1000; i++ {
		s1 := NewSized(10)
		s2 := NewSized(10)
		s1.Set(1)
		s1.Set(2)
		s1.Set(3)

		s2.Set(2)
		s2.Set(3)
		s2.Set(4)

		s := Union([]intSet{s1, s2})
		expect.Expect(s.Exists(1)).To.Equal(true)
		expect.Expect(s.Exists(2)).To.Equal(true)
		expect.Expect(s.Exists(3)).To.Equal(true)
		expect.Expect(s.Exists(4)).To.Equal(true)
		expect.Expect(s.Exists(5)).To.Equal(false)
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
		Intersect([]intSet{s1, s2})
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
