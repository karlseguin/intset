// Package intset provides a specialized set for integers or runes
package intset

const (
	defaultBucketSize   int = 4
	defaultBucketGrowBy int = 1
)

// Config defines the configuration for creating a new set
//
// Smaller values for bucketSize will speed up lookups but
// also increase memory usage. Smaller values for bucketGrowBy will slow down
// the set capacity growth rate but also slow down insertions.
type Config struct {
	bucketSize   int
	bucketGrowBy int
}

// NewConfig creates a new config with usable defaults
func NewConfig() *Config {
	return &Config{bucketSize: defaultBucketSize, bucketGrowBy: defaultBucketGrowBy}
}

// BucketSize sets the initial bucket size
func (c *Config) BucketSize(size uint32) *Config {
	if size == 0 { // bucketSize must be positive int
		c.bucketSize = defaultBucketSize
	} else {
		c.bucketSize = int(size)
	}
	return c
}

// BucketGrowBy sets the amount a bucket will grow by when full
func (c *Config) BucketGrowBy(growBy uint32) *Config {
	c.bucketGrowBy = int(growBy)
	return c
}

// Default is a default Config which favors probing performance
// at the cost of memory.
var Default = NewConfig()
