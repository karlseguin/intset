// Package intset provides a specialized set for integers or runes
package intset

import "sort"

// SetRune defines rune set methods
type SetRune interface {
	Len() int
	Exists(value rune) bool
	Each(f func(value rune))
}

// SetsRune is array of Rune
type SetsRune []SetRune

func (s SetsRune) Len() int {
	return len(s)
}

func (s SetsRune) Less(i, j int) bool {
	return s[i].Len() < s[j].Len()
}

func (s SetsRune) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Rune stores rune set data
type Rune struct {
	mask    rune
	buckets [][]rune
	length  int
	growBy  int
}

// NewRune creates an empty rune set with target capacity specified by size using default configuration
func NewRune(size rune) *Rune {
	return NewRuneConfig(size, Default)
}

// NewRuneConfig creates an empty rune set with target capacity specified by size
func NewRuneConfig(size rune, config *Config) *Rune {
	if size < rune(config.bucketSize) {
		size = rune(config.bucketSize)
	}
	count := upTwo(int(size) / config.bucketSize)
	s := &Rune{
		mask:    rune(count) - 1,
		buckets: make([][]rune, count),
		growBy:  config.bucketGrowBy,
	}
	return s
}

// Set adds a value to the rune set
func (s *Rune) Set(value rune) {
	index := value & s.mask
	bucket := s.buckets[index]
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	l := len(bucket)
	if cap(bucket) == l {
		n := make([]rune, l, l+s.growBy)
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

// Remove returns true if the value existed in the rune set before being removed
func (s *Rune) Remove(value rune) bool {
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
func (s *Rune) Exists(value rune) bool {
	return s.exists(value, s.buckets[value&s.mask])
}

// Len returns the total number of elements in the set
func (s Rune) Len() int {
	return s.length
}

// Each iterates through the set items and applies function f to each set item
func (s Rune) Each(f func(value rune)) {
	for _, bucket := range s.buckets {
		for _, value := range bucket {
			f(value)
		}
	}
}

func (s Rune) index(value rune, bucket []rune) (int, bool) {
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

func (s Rune) exists(value rune, bucket []rune) bool {
	l := len(bucket)
	if l == 0 {
		return false
	}

	l = l / 2
	v := bucket[l]
	if value == v {
		return true
	}
	if value > v {
		bucket = bucket[l:]
	}

	for _, v = range bucket {
		if v >= value {
			return v == value
		}
	}
	return false
}

// IntersectRune returns the intersection of an array of sets
func IntersectRune(sets SetsRune) *Rune {
	sort.Sort(sets)
	a, l := sets[0], sets.Len()
	values := make([]rune, 0, a.Len())
	a.Each(func(value rune) {
		for i := 1; i < l; i++ {
			if sets[i].Exists(value) == false {
				return
			}
		}
		values = append(values, value)
	})
	s := NewRune(rune(len(values)))
	for _, value := range values {
		s.Set(value)
	}
	return s
}

// UnionRune returns the union of an array of sets
func UnionRune(sets SetsRune) *Rune {
	values := make(map[rune]struct{}, sets[0].Len())
	for i := 0; i < sets.Len(); i++ {
		sets[i].Each(func(value rune) {
			values[value] = struct{}{}
		})
	}
	s := NewRune(rune(len(values)))
	for value := range values {
		s.Set(value)
	}
	return s
}
