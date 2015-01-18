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
	s := NewSized(SIZED_BUCKET_SIZE - 1)
	s.Set(32)
	Expect(s.Exists(32)).To.Equal(true)
	Expect(s.Exists(33)).To.Equal(false)
}

func Benchmark_SizedPopulate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSized(1000000)
		s.Set(i)
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
		if s.Exists(rand.Intn(1000000)) == false {
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
		if s.Exists(rand.Intn(1000000)) == false {
			misses++
		}
	}
}
