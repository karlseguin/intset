// Package intset provides a specialized set for integers or runes
package intset

// BucketConfig specifies optional parameters for tuning intset performance.
//
// Smaller values for bucketSize and bucketMultiplier will speed up lookups but
// also increase memory usage. Smaller values for bucketGrowBy will slow down
// the set capacity growth rate but also slow down insertions.
//
// Any omitted parameters will be set to default values by setDefaultBucketConfig()
type BucketConfig struct {
	bucketSize       int
	bucketMultiplier int
	bucketGrowBy     int
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

func setDefaultBucketConfig(bucketSize int, bucketMultiplier int, bucketGrowBy int) (int, int, int) {
	// Default values for a reasonable space-time tradeoff
	if bucketSize < 1 {
		bucketSize = 4
	}
	if bucketMultiplier < 1 {
		bucketMultiplier = 1
	}
	if bucketGrowBy < 1 {
		bucketGrowBy = 1
	}
	return bucketSize, bucketMultiplier, bucketGrowBy
}
