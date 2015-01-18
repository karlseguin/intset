// integer set
package intset

type Sized32 struct {
	mask    uint32
	buckets [][]uint32
}

func NewSized32(size uint32) *Sized32 {
	if size < SIZED_BUCKET_SIZE {
		//random, no clue what to make it
		size = uint32(SIZED_BUCKET_SIZE * 2)
	}
	count := upTwo(int(size) / SIZED_BUCKET_SIZE)
	s := &Sized32{
		mask:    uint32(count) - 1,
		buckets: make([][]uint32, count),
	}
	return s
}

func (s *Sized32) Set(value uint32) {
	index := value & s.mask
	bucket := s.buckets[index]
	l := len(bucket)
	if l == 0 {
		s.buckets[index] = []uint32{value}
		return
	}
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
	arr := make([]uint32, l+1)
	copy(arr, bucket[:position])
	arr[position] = value
	copy(arr[position+1:], bucket[position:])
	s.buckets[index] = arr
}

func (s *Sized32) Exists(value uint32) bool {
	bucket := s.buckets[value&s.mask]
	_, exists := s.index(value, bucket)
	return exists
}

func (s Sized32) index(value uint32, bucket []uint32) (int, bool) {
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
