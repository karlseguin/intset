// integer set
package intset

const (
	SIZED_BUCKET_SIZE = 32
)

type Sized struct {
	mask    int
	buckets [][]int
}

func NewSized(size int) *Sized {
	if size < SIZED_BUCKET_SIZE {
		//random, no clue what to make it
		size = SIZED_BUCKET_SIZE * 2
	}
	count := upTwo(size / SIZED_BUCKET_SIZE)
	s := &Sized{
		mask:    count - 1,
		buckets: make([][]int, count),
	}
	return s
}

func (s *Sized) Set(value int) {
	index := value & s.mask
	bucket := s.buckets[index]
	l := len(bucket)
	if l == 0 {
		s.buckets[index] = []int{value}
		return
	}
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	arr := make([]int, l+1)
	copy(arr, bucket[:position])
	arr[position] = value
	copy(arr[position+1:], bucket[position:])
	s.buckets[index] = arr
}

func (s *Sized) Exists(value int) bool {
	bucket := s.buckets[value&s.mask]
	_, exists := s.index(value, bucket)
	return exists
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
