package intset

import (
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
	s := NewRune(rune(BUCKET_SIZE) - 1)
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

	s := IntersectRune([]SetRune{s1, s2})
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

	s := UnionRune([]SetRune{s1, s2})
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
	misses := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Exists(rune(i%1000000)) == false {
			misses++
		}
	}
}

func Benchmark_RuneSparseExists(b *testing.B) {
	s := NewRune(1000000)
	for i := rune(0); i < 1000000; i++ {
		if i%10 == 0 {
			s.Set(i)
		}
	}
	misses := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Exists(rune(i%1000000)) == false {
			misses++
		}
	}
}
