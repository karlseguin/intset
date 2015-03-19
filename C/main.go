package main

/*
#include "intset.h"
*/
import "C"

import (
	"fmt"
	"github.com/karlseguin/intset"
	"runtime"
	"testing"
)

type Set interface {
	Set(value int)
	Exists(value int) bool
}

type MapSet map[int]struct{}

func NewMap(count int) MapSet {
	return make(map[int]struct{}, count)
}

func (s MapSet) Set(value int) {
	s[value] = struct{}{}
}

func (s MapSet) Exists(value int) bool {
	_, exists := s[value]
	return exists
}

func (s1 MapSet) Intersect(s2 MapSet) MapSet {
	count := 0
	values := make([]int, len(s1))
	for k, _ := range s1 {
		if _, exists := s2[k]; exists {
			values[count] = k
		}
	}

	n := NewMap(count)
	for i := 0; i < count; i++ {
		n[values[i]] = struct{}{}
	}
	return n
}

type IntSet C.IntSet

func NewIntSet(size uint) *IntSet {
	return (*IntSet)(C.intset_new(C.longlong(size)))
}

func (s *IntSet) Set(value int) {
	C.intset_set((*C.IntSet)(s), C.longlong(value))
}

func (s *IntSet) Len() int {
	return int(s.length)
}

func (s *IntSet) Exists(value int) bool {
	return C.intset_exists((*C.IntSet)(s), C.longlong(value)) == 1
}

func (s1 *IntSet) Intersect(s2 *IntSet) *IntSet {
	return (*IntSet)(C.intset_intersect((*C.IntSet)(s1), (*C.IntSet)(s2)))
}

func (s *IntSet) Close() {
	C.intset_free((*C.IntSet)(s))
}

func main() {
	benchmark()
}

func mem() {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	fmt.Println(stats.HeapAlloc)
}

func benchmark() {
	fmt.Println("go populate", testing.Benchmark(SizedPopulate(intset.NewSized(10000000))))
	fmt.Println("go dense", testing.Benchmark(SizedDenseExists(intset.NewSized(10000000))))
	fmt.Println("go sparse", testing.Benchmark(SizedSparseExists(intset.NewSized(10000000))))
	fmt.Println("go intersect", testing.Benchmark(GoIntersect(intset.NewSized(10000000), intset.NewSized(10000000))))

	fmt.Println("map populate", testing.Benchmark(SizedPopulate(NewMap(10000000))))
	fmt.Println("map dense", testing.Benchmark(SizedDenseExists(NewMap(10000000))))
	fmt.Println("map sparse", testing.Benchmark(SizedSparseExists(NewMap(10000000))))
	fmt.Println("map intersect", testing.Benchmark(MapIntersect(NewMap(10000000), NewMap(10000000))))

	s := NewIntSet(10000000)
	fmt.Println("c populate", testing.Benchmark(SizedPopulate(s)))
	s.Close()

	s = NewIntSet(10000000)
	fmt.Println("c dense", testing.Benchmark(SizedDenseExists(s)))
	s.Close()

	s = NewIntSet(10000000)
	fmt.Println("c sparse", testing.Benchmark(SizedSparseExists(s)))
	s.Close()

	s1 := NewIntSet(10000000)
	s2 := NewIntSet(10000000)
	fmt.Println("c intersect", testing.Benchmark(CIntersect(s1, s2)))
	s1.Close()
	s2.Close()
}

func SizedPopulate(s Set) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Set(i % 10000000)
		}
	}
}

func SizedDenseExists(s Set) func(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		s.Set(i)
	}

	return func(b *testing.B) {
		b.ResetTimer()
		misses := 0
		for i := 0; i < b.N; i++ {
			if s.Exists(i%1000000) == false {
				misses++
			}
		}
	}
}

func SizedSparseExists(s Set) func(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		if i%10 == 0 {
			s.Set(i)
		}
	}

	return func(b *testing.B) {
		b.ResetTimer()
		misses := 0
		for i := 0; i < b.N; i++ {
			if s.Exists(i%1000000) == false {
				misses++
			}
		}
	}
}

func GoIntersect(s1 *intset.Sized, s2 *intset.Sized) func(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		if i%3 == 0 {
			s1.Set(i)
		}
		if i%4 == 0 {
			s2.Set(i)
		}
	}

	return func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s1.Intersect(s2)
		}
	}
}

func CIntersect(s1 *IntSet, s2 *IntSet) func(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		if i%3 == 0 {
			s1.Set(i)
		}
		if i%4 == 0 {
			s2.Set(i)
		}
	}

	return func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s1.Intersect(s2)
		}
	}
}

func MapIntersect(s1 MapSet, s2 MapSet) func(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		if i%3 == 0 {
			s1.Set(i)
		}
		if i%4 == 0 {
			s2.Set(i)
		}
	}

	return func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s1.Intersect(s2)
		}
	}
}
