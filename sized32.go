// integer set
package intset

import (
	"sort"
)

type Set32 interface {
	Len() int
	Exists(value uint32) bool
	Each(f func(value uint32))
}

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

type Sized32 struct {
	mask    uint32
	buckets [][]uint32
	length  int
}

func NewSized32(size uint32) *Sized32 {
	if size < uint32(BUCKET_SIZE) {
		//random, no clue what to make it
		size = uint32(BUCKET_SIZE * 2)
	}
	count := upTwo(int(size) / BUCKET_SIZE)
	s := &Sized32{
		mask:    uint32(count) - 1,
		buckets: make([][]uint32, count),
	}
	return s
}

func (s *Sized32) Set(value uint32) {
	index := value & s.mask
	bucket := s.buckets[index]
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	l := len(bucket)
	if cap(bucket) == l {
		n := make([]uint32, l, l+BUCKET_GROW_BY)
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

// returns true if the value existed
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

func (s *Sized32) Exists(value uint32) bool {
	bucket := s.buckets[value&s.mask]
	_, exists := s.index(value, bucket)
	return exists
}

func (s Sized32) Len() int {
	return s.length
}

// Iterate through the set items
func (s Sized32) Each(f func(value uint32)) {
	for i, li := 0, len(s.buckets); i < li; i++ {
		bucket := s.buckets[i]
		for j, lj := 0, len(bucket); j < lj; j++ {
			f(bucket[j])
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

	offset, i := 0, 0
	if value > v {
		offset = l
		bucket = bucket[offset:]
	}

	for i, v = range bucket {
		if v < value {
			continue
		}
		return offset + i, v == value
	}
	return offset + i + 1, false
}

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
