// integer set
package intset

var (
	BUCKET_SIZE    = 32
	BUCKET_GROW_BY = 8
)

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
