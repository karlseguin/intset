package intset

import (
	. "github.com/karlseguin/expect"
	"math/rand"
	"testing"
)

type SizedTest struct{}

func Test_Sized(t *testing.T) {
	Expectify(new(SizedTest), t)
}

func (_ SizedTest) SetsAValue() {
	s := NewSized(20)
	for i := 0; i < 30; i++ {
		s.Set(i)
		Expect(s.Exists(i)).To.Equal(true)
	}
	for i := 0; i < 30; i++ {
		Expect(s.Exists(i)).To.Equal(true)
	}
}

func (_ SizedTest) Exists() {
	s := NewSized(20)
	for i := 0; i < 10; i++ {
		Expect(s.Exists(i)).To.Equal(false)
		s.Set(i)
	}
	for i := 0; i < 10; i++ {
		Expect(s.Exists(i)).To.Equal(true)
	}
}

func (_ SizedTest) SizeLessThanBucket() {
	s := NewSized(BUCKET_SIZE - 1)
	s.Set(32)
	Expect(s.Exists(32)).To.Equal(true)
	Expect(s.Exists(33)).To.Equal(false)
}

func (_ SizedTest) RemoveNonMembers() {
	s := NewSized(100)
	Expect(s.Remove(329)).To.Equal(false)
}

func (_ SizedTest) RemovesMembers() {
	s := NewSized(100)
	for i := 0; i < 10; i++ {
		s.Set(i)
	}
	Expect(s.Remove(20)).To.Equal(false)
	Expect(s.Remove(2)).To.Equal(true)
	Expect(s.Remove(2)).To.Equal(false)
	Expect(s.Exists(2)).To.Equal(false)
	Expect(s.Len()).To.Equal(9)
}

func (_ SizedTest) IntersectsTwoSets() {
	s1 := NewSized(10)
	s2 := NewSized(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Intersect([]Set{s1, s2})
	Expect(s.Exists(1)).To.Equal(false)
	Expect(s.Exists(2)).To.Equal(true)
	Expect(s.Exists(3)).To.Equal(true)
	Expect(s.Exists(4)).To.Equal(false)
	Expect(s.Exists(5)).To.Equal(false)
}

func (_ SizedTest) UnionsTwoSets() {
	s1 := NewSized(10)
	s2 := NewSized(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Union([]Set{s1, s2})
	Expect(s.Exists(1)).To.Equal(true)
	Expect(s.Exists(2)).To.Equal(true)
	Expect(s.Exists(3)).To.Equal(true)
	Expect(s.Exists(4)).To.Equal(true)
	Expect(s.Exists(5)).To.Equal(false)
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
	misses := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Exists(i%1000000) == false {
			misses++
		}
	}
}

func Benchmark_SizedSparseExists(b *testing.B) {
	s := NewSized(1000000)
	for i := 0; i < 1000000; i++ {
		if i%10 == 0 {
			s.Set(i)
		}
	}
	misses := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Exists(i%1000000) == false {
			misses++
		}
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
