# IntSet

[![Go Reference](https://img.shields.io/badge/go-reference-blue?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/karlseguin/intset)
[![Go Report Card](https://goreportcard.com/badge/github.com/karlseguin/intset?style=for-the-badge)](https://goreportcard.com/report/github.com/karlseguin/intset)
[![GitHub license](https://img.shields.io/badge/LICENSE-MIT-GREEN?style=for-the-badge)](license.txt)

A specialized set for integers or runes, ideal when:

- The number of elements is known ahead of time (or a good approximation)
- The number of elements doesn't change drastically over time
- The values are naturally random

As long as the number of elements within the set stays close to the originally specified size (I don't know the magic number, so let's say ±10%), and that they stay evenly distributed. the set will exhibit good read and write performance, as well as decent memory usage. When packed, read performance is roughly 7 times better than a map[int]struct{}.

```go
set := intset.NewSized(1000000, BucketConfig{})  // or intset.NewSized32(1000000, BucketConfig{}) or intset.NewRune(1000000, BucketConfig{})
set.Set(32)
set.Exists(32)
```

## Methods

The `int`, `uint32` and `rune` variations have the same API.

- `Set(int)` or `Set(uint32)` or `Set(rune)`
- `Exists(int) bool` or `Exists(uint32) bool` or `Exists(rune) bool`
- `Remove(int) bool` or `Remove(uint32) bool` or `Remove(rune) bool`
- `Len() int`
- `Each(f func(value int))` or `Each(f func(value uint32))` or `Each(f func(value rune))`

## Intersections and Unions

Two or more sets can be intersected by calling `Intersect`, `Intersect32`, or `IntersectRune`. This is largely a reference implementation and callers should consider implementing their own. For example, maybe you want to stop after finding X matches, want to use a pooled array object to hold intermediary objects, or are fine with getting an array back (rather than a set) (all of which should result in much better performance).

The method is called via:

```go
result := intset.Intersect([]Set{s1, s2})
// or
result := intset.Intersect32([]Set32{s1, s2})
// or
result := intset.IntersectRune([]Set32{s1, s2})
```

`Union`, `Union32`, and `UnionRune` can be similarly used.
