package intset

import (
	"testing"

	expect "github.com/karlseguin/expect"
)

type Sized32Test struct{}

func Test_Sized32(t *testing.T) {
	expect.Expectify(new(Sized32Test), t)
}

func (Sized32Test) SetsAValue() {
	s := NewSized32(20)
	for i := uint32(0); i < 30; i++ {
		s.Set(i)
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
	for i := uint32(0); i < 30; i++ {
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
}

func (Sized32Test) Exists() {
	s := NewSized32(20)
	for i := uint32(0); i < 10; i++ {
		expect.Expect(s.Exists(i)).To.Equal(false)
		s.Set(i)
	}
	for i := uint32(0); i < 10; i++ {
		expect.Expect(s.Exists(i)).To.Equal(true)
	}
}

func (Sized32Test) SizeLessThanBucket() {
	s := NewSized32(uint32(bucketSize) - 1)
	s.Set(32)
	expect.Expect(s.Exists(32)).To.Equal(true)
	expect.Expect(s.Exists(33)).To.Equal(false)
}

func (Sized32Test) RemoveNonMembers() {
	s := NewSized32(100)
	expect.Expect(s.Remove(329)).To.Equal(false)
}

func (Sized32Test) RemovesMembers() {
	s := NewSized32(100)
	for i := uint32(0); i < 10; i++ {
		s.Set(i)
	}
	expect.Expect(s.Remove(20)).To.Equal(false)
	expect.Expect(s.Remove(2)).To.Equal(true)
	expect.Expect(s.Remove(2)).To.Equal(false)
	expect.Expect(s.Exists(2)).To.Equal(false)
	expect.Expect(s.Len()).To.Equal(9)
}

func (Sized32Test) IntersectsTwoSets() {
	s1 := NewSized32(10)
	s2 := NewSized32(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Intersect32([]Set32{s1, s2})
	expect.Expect(s.Exists(1)).To.Equal(false)
	expect.Expect(s.Exists(2)).To.Equal(true)
	expect.Expect(s.Exists(3)).To.Equal(true)
	expect.Expect(s.Exists(4)).To.Equal(false)
	expect.Expect(s.Exists(5)).To.Equal(false)
}

func (Sized32Test) UnionsTwoSets() {
	s1 := NewSized32(10)
	s2 := NewSized32(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Union32([]Set32{s1, s2})
	expect.Expect(s.Exists(1)).To.Equal(true)
	expect.Expect(s.Exists(2)).To.Equal(true)
	expect.Expect(s.Exists(3)).To.Equal(true)
	expect.Expect(s.Exists(4)).To.Equal(true)
	expect.Expect(s.Exists(5)).To.Equal(false)
}

func Benchmark_Sized32Populate(b *testing.B) {
	s := NewSized32(10000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(uint32(i % 10000000))
	}
}

func Benchmark_Sized32DenseExists(b *testing.B) {
	s := NewSized32(1000000)
	for i := uint32(0); i < 1000000; i++ {
		s.Set(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exists(uint32(i % 1000000))
	}
}

func Benchmark_Sized32SparseExists(b *testing.B) {
	s := NewSized32(1000000)
	for i := uint32(0); i < 1000000; i++ {
		if i%10 == 0 {
			s.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exists(uint32(i % 1000000))
	}
}
