package intset

import (
	. "github.com/karlseguin/expect"
	"math/rand"
	"testing"
)

type Sized32Test struct{}

func Test_Sized32(t *testing.T) {
	Expectify(new(Sized32Test), t)
}

func (_ Sized32Test) SetsAValue() {
	s := NewSized32(20)
	for i := uint32(0); i < 30; i++ {
		s.Set(i)
		Expect(s.Exists(i)).To.Equal(true)
	}
	for i := uint32(0); i < 30; i++ {
		Expect(s.Exists(i)).To.Equal(true)
	}
}

func (_ Sized32Test) Exists() {
	s := NewSized32(20)
	for i := uint32(0); i < 10; i++ {
		Expect(s.Exists(i)).To.Equal(false)
		s.Set(i)
	}
	for i := uint32(0); i < 10; i++ {
		Expect(s.Exists(i)).To.Equal(true)
	}
}

func (_ Sized32Test) SizeLessThanBucket() {
	s := NewSized32(uint32(BUCKET_SIZE) - 1)
	s.Set(32)
	Expect(s.Exists(32)).To.Equal(true)
	Expect(s.Exists(33)).To.Equal(false)
}

func Benchmark_Sized32Populate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSized32(1000000)
		s.Set(uint32(i))
	}
}

func Benchmark_Sized32DenseExists(b *testing.B) {
	s := NewSized32(1000000)
	for i := uint32(0); i < 1000000; i++ {
		s.Set(i)
	}
	misses := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Exists(uint32(rand.Int31n(1000000))) == false {
			misses++
		}
	}
}

func Benchmark_Sized32SparseExists(b *testing.B) {
	s := NewSized32(1000000)
	for i := uint32(0); i < 1000000; i++ {
		if i%10 == 0 {
			s.Set(i)
		}
	}
	misses := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Exists(uint32(rand.Int31n(1000000))) == false {
			misses++
		}
	}
}
