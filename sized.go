// Package intset provides a specialized set for integers or runes
package intset

import "sort"

// Set defines int set methods
type Set interface {
	Len() int
	Exists(value int) bool
	Each(f func(value int))
}

// Sets is array of Set
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

// Sized stores int set data
type Sized struct {
	mask    int
	buckets [][]int
	length  int
	growBy  int
}

// NewSized creates an empty int set with target capacity specified by size using default configuration
func NewSized(size int) *Sized {
	return NewSizedConfig(size, Default)
}

// NewSizedConfig creates an empty int set with target capacity specified by size
func NewSizedConfig(size int, config *Config) *Sized {
	if size < config.bucketSize {
		size = config.bucketSize
	}
	count := upTwo(size / config.bucketSize)
	s := &Sized{
		mask:    count - 1,
		buckets: make([][]int, count),
		growBy:  config.bucketGrowBy,
	}
	return s
}

// Set adds a value to the int set
func (s *Sized) Set(value int) {
	index := value & s.mask
	bucket := s.buckets[index]
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	l := len(bucket)
	if cap(bucket) == l {
		n := make([]int, l, l+s.growBy)
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

// Remove returns true if the value existed in the int set before being removed
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

// Exists returns true if the value exists in the set
func (s *Sized) Exists(value int) bool {
	return s.exists(value, s.buckets[value&s.mask])
}

// Len returns the total number of elements in the set
func (s Sized) Len() int {
	return s.length
}

// Each iterates through the set items and applies function f to each set item
func (s Sized) Each(f func(value int)) {
	for _, bucket := range s.buckets {
		for _, value := range bucket {
			f(value)
		}
	}
}

func (s Sized) index(value int, bucket []int) (int, bool) {
	l := len(bucket)
	if l == 0 {
		return 0, false
	}

	l = l / 2
	v := bucket[l]
	if value == v {
		return l, true
	}
	if value > v {
		bucket = bucket[l:]
	} else {
		l = 0
	}

	var i int
	for i, v = range bucket {
		if v >= value {
			return l + i, v == value
		}
	}
	return l + i + 1, false
}

func (s Sized) exists(value int, bucket []int) bool {
	if len(bucket) == 0 {
		return false
	}

	for _, v := range bucket {
		if v >= value {
			return v == value
		}
	}
	return false
}

// Intersect returns the intersection of an array of sets
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

// Union returns the union of an array of sets
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
