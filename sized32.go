// integer set
package intset

type Sized32 struct {
	mask    uint32
	buckets [][]uint32
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
	l := len(bucket)
	if l == 0 {
		s.buckets[index] = []uint32{value}
		return
	}
	position, exists := s.index(value, bucket)
	if exists {
		return
	}
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
	s.buckets[index] = bucket
}

func (s *Sized32) Exists(value uint32) bool {
	bucket := s.buckets[value&s.mask]
	_, exists := s.index(value, bucket)
	return exists
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
	return true
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
