// integer set
package intset

import (
	"sort"
)

var (
	BUCKET_SIZE    = 32
	BUCKET_GROW_BY = 8
)

type Set interface {
	Len() int
	Exists(value int) bool
	Each(f func(value int))
}

type Sets []Set

func (s Sets) Len() int {
	return len(s)
}

func (s Sets) Less(i, j int) bool {
	return s[i].Len() < s[j].Len()
}

func (s Sets) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Sized struct {
	mask    int
	buckets [][]int
	length  int
}

func NewSized(size int) *Sized {
	if size < BUCKET_SIZE {
		//random, no clue what to make it
		size = BUCKET_SIZE * 2
	}
	count := upTwo(size / BUCKET_SIZE)
	s := &Sized{
		mask:    count - 1,
		buckets: make([][]int, count),
	}
	return s
}

// Sets a value
func (s *Sized) Set(value int) {
	index := value & s.mask
	bucket := s.buckets[index]
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	l := len(bucket)
	if cap(bucket) == l {
		n := make([]int, l, l+BUCKET_GROW_BY)
		copy(n, bucket)
		bucket = n
	}
	bucket = append(bucket, value)
	if position != l {
		copy(bucket[position+1:], bucket[position:])
		bucket[position] = value
	}
	s.length++
	s.buckets[index] = bucket
}

// Returns true if the value existed before being removed
func (s *Sized) Remove(value int) bool {
	index := value & s.mask
	bucket := s.buckets[index]
	position, exists := s.index(value, bucket)
	if exists == false {
		return false
	}
	l := len(bucket) - 1
	bucket[position], bucket[l] = bucket[l], bucket[position]
	s.buckets[index] = bucket[:l]
	s.length--
	return true
}

// Returns true if the value exists
func (s *Sized) Exists(value int) bool {
	bucket := s.buckets[value&s.mask]
	_, exists := s.index(value, bucket)
	return exists
}

// Total number of elements in the set
func (s Sized) Len() int {
	return s.length
}

// Iterate through the set items
func (s Sized) Each(f func(value int)) {
	for i, li := 0, len(s.buckets); i < li; i++ {
		bucket := s.buckets[i]
		for j, lj := 0, len(bucket); j < lj; j++ {
			f(bucket[j])
		}
	}
}

func (s Sized) index(value int, bucket []int) (int, bool) {
	l := len(bucket)
	for i := 0; i < l; i++ {
		v := bucket[i]
		if v == value {
			return i, true
		}
		if v > value {
			return i, false
		}
	}
	return l, false
}

func Intersect(sets Sets) *Sized {
	sort.Sort(sets)
	a, l := sets[0], sets.Len()
	values := make([]int, 0, a.Len())
	a.Each(func(value int) {
		for i := 1; i < l; i++ {
			if sets[i].Exists(value) == false {
				return
			}
		}
		values = append(values, value)
	})
	s := NewSized(len(values))
	for _, value := range values {
		s.Set(value)
	}
	return s
}

func Union(sets Sets) *Sized {
	values := make(map[int]struct{}, sets[0].Len())
	for i := 0; i < sets.Len(); i++ {
		sets[i].Each(func(value int) {
			values[value] = struct{}{}
		})
	}
	s := NewSized(len(values))
	for value := range values {
		s.Set(value)
	}
	return s
}

// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func upTwo(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
