// integer set
package intset

import "sort"

type SetRune interface {
	Len() int
	Exists(value rune) bool
	Each(f func(value rune))
}

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

type Rune struct {
	mask    rune
	buckets [][]rune
	length  int
}

func NewRune(size rune) *Rune {
	if size < rune(bucketSize) {
		//random, no clue what to make it
		size = rune(bucketSize) * rune(bucketMultiplier)
	}
	count := upTwo(int(size) / bucketSize)
	s := &Rune{
		mask:    rune(count) - 1,
		buckets: make([][]rune, count),
	}
	return s
}

func (s *Rune) Set(value rune) {
	index := value & s.mask
	bucket := s.buckets[index]
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	l := len(bucket)
	if cap(bucket) == l {
		n := make([]rune, l, l+bucketGrowBy)
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

func (s *Rune) Exists(value rune) bool {
	return s.exists(value, s.buckets[value&s.mask])
}

func (s Rune) Len() int {
	return s.length
}

// Iterate through the set items
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
