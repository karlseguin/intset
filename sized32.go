// Package intset provides a specialized set for integers or runes
package intset

import "sort"

// Set32 defines uint32 set methods
type Set32 interface {
	Len() int
	Exists(value uint32) bool
	Each(f func(value uint32))
}

// Sets32 is array of Set32
type Sets32 []Set32

func (s Sets32) Len() int {
	return len(s)
}

func (s Sets32) Less(i, j int) bool {
	return s[i].Len() < s[j].Len()
}

func (s Sets32) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sized32 stores uint32 set data
type Sized32 struct {
	mask    uint32
	buckets [][]uint32
	length  int
	growBy  int
}

// NewSized32 creates an empty int set with target capacity specified by size using default configuration
func NewSized32(size uint32) *Sized32 {
	return NewSized32Config(size, Default)
}

// NewSized32Config creates an empty uint32 set with target capacity specified by size
func NewSized32Config(size uint32, config *Config) *Sized32 {
	if size < uint32(config.bucketSize) {
		size = uint32(config.bucketSize)
	}
	count := upTwo(int(size) / config.bucketSize)
	s := &Sized32{
		mask:    uint32(count) - 1,
		buckets: make([][]uint32, count),
		growBy:  config.bucketGrowBy,
	}
	return s
}

// Set adds a value to the uint32 set
func (s *Sized32) Set(value uint32) {
	index := value & s.mask
	bucket := s.buckets[index]
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	l := len(bucket)
	if cap(bucket) == l {
		n := make([]uint32, l, l+s.growBy)
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

// Remove returns true if the value existed in the uint32 set before being removed
func (s *Sized32) Remove(value uint32) bool {
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
func (s *Sized32) Exists(value uint32) bool {
	return s.exists(value, s.buckets[value&s.mask])
}

// Len returns the total number of elements in the set
func (s Sized32) Len() int {
	return s.length
}

// Each iterates through the set items and applies function f to each set item
func (s Sized32) Each(f func(value uint32)) {
	for _, bucket := range s.buckets {
		for _, value := range bucket {
			f(value)
		}
	}
}

func (s Sized32) index(value uint32, bucket []uint32) (int, bool) {
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

func (s Sized32) exists(value uint32, bucket []uint32) bool {
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

// Intersect32 returns the intersection of an array of sets
func Intersect32(sets Sets32) *Sized32 {
	sort.Sort(sets)
	a, l := sets[0], sets.Len()
	values := make([]uint32, 0, a.Len())
	a.Each(func(value uint32) {
		for i := 1; i < l; i++ {
			if sets[i].Exists(value) == false {
				return
			}
		}
		values = append(values, value)
	})
	s := NewSized32(uint32(len(values)))
	for _, value := range values {
		s.Set(value)
	}
	return s
}

// Union32 returns the union of an array of sets
func Union32(sets Sets32) *Sized32 {
	values := make(map[uint32]struct{}, sets[0].Len())
	for i := 0; i < sets.Len(); i++ {
		sets[i].Each(func(value uint32) {
			values[value] = struct{}{}
		})
	}
	s := NewSized32(uint32(len(values)))
	for value := range values {
		s.Set(value)
	}
	return s
}
