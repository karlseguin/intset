// Package intset provides a specialized set for integers or runes
package intset

const bucketSize, bucketGrowBy, bucketMultiplier int = 4, 1, 1

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
