package intset

import (
	"math/rand"
	"testing"
)

func Test_Sized32_SetsAValue(t *testing.T) {
	s := NewSized32(20)
	for i := uint32(0); i < 30; i++ {
		s.Set(i)
		AssertTrue(t, s.Exists(i))
	}
	for i := uint32(0); i < 30; i++ {
		AssertTrue(t, s.Exists(i))
	}
}

func Test_Sized32_Exists(t *testing.T) {
	s := NewSized32Config(20, NewConfig())
	for i := uint32(0); i < 10; i++ {
		AssertFalse(t, s.Exists(i))
		s.Set(i)
	}
	for i := uint32(0); i < 10; i++ {
		AssertTrue(t, s.Exists(i))
	}
}

func Test_Sized32_SizeLessThanBucket(t *testing.T) {
	config := NewConfig().BucketSize(8).BucketGrowBy(4)
	s := NewSized32Config(7, config)
	s.Set(32)
	AssertTrue(t, s.Exists(32))
	AssertFalse(t, s.Exists(33))
}

func Test_Sized32_RemoveNonMembers(t *testing.T) {
	s := NewSized32(100)
	AssertFalse(t, s.Remove(329))
}

func Test_Sized32_RemovesMembers(t *testing.T) {
	s := NewSized32(100)
	for i := uint32(0); i < 10; i++ {
		s.Set(i)
	}
	AssertFalse(t, s.Remove(20))
	AssertTrue(t, s.Remove(2))
	AssertFalse(t, s.Remove(2))
	AssertFalse(t, s.Exists(2))
	AssertEqual(t, s.Len(), 9)
}

func Test_Sized32_IntersectsTwoSets(t *testing.T) {
	s1 := NewSized32(10)
	s2 := NewSized32(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Intersect32([]Set32{s1, s2})
	AssertFalse(t, s.Exists(1))
	AssertTrue(t, s.Exists(2))
	AssertTrue(t, s.Exists(3))
	AssertFalse(t, s.Exists(4))
	AssertFalse(t, s.Exists(5))
}

func Test_Sized32_UnionsTwoSets(t *testing.T) {
	s1 := NewSized32(10)
	s2 := NewSized32(10)
	s1.Set(1)
	s1.Set(2)
	s1.Set(3)

	s2.Set(2)
	s2.Set(3)
	s2.Set(4)

	s := Union32([]Set32{s1, s2})
	AssertTrue(t, s.Exists(1))
	AssertTrue(t, s.Exists(2))
	AssertTrue(t, s.Exists(3))
	AssertTrue(t, s.Exists(4))
	AssertFalse(t, s.Exists(5))
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

func Benchmark_Sized32DenseIntersect(b *testing.B) {
	s1 := NewSized32(100000)
	for i := uint32(0); i < 100000; i++ {
		if rand.Intn(10) != 0 {
			s1.Set(i)
		}
	}
	s2 := NewSized32(1000)
	for i := uint32(0); i < 1000; i++ {
		if rand.Intn(10) != 0 {
			s2.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Intersect32([]Set32{s1, s2})
	}
}

// Benchmarks for map[uint32]struct{}

func Benchmark_Sized32MapDenseExists(b *testing.B) {
	s := make(map[uint32]struct{})
	for i := uint32(0); i < 1000000; i++ {
		s[i] = struct{}{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s[uint32(i)]
	}
}

func Benchmark_Sized32MapSparseExists(b *testing.B) {
	s := make(map[uint32]struct{})
	for i := uint32(0); i < 1000000; i++ {
		if i%10 == 0 {
			s[i] = struct{}{}
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s[uint32(i%1000000)]
	}
}
