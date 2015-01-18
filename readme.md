A specialized set for integers, ideal when:

- The number of elements is known ahead of time (or a good approximation)
- The number of elements doesn't change drastically over time
- The values are naturally random

As long as the number of elements within the set stays close to the originally specified size (I don't know the magic number, so let's say ±10%), and that they stay evenly distributed. the set will exhibit good read and write performance, as well as decent memory usage. When packed, read performance is roughly 2x better than a map[int]struct{} with memory usage less than 1/2.


```go
set := intset.Sized(1000000)  //or inteet.Sized32(100000)
set.Set(32)
set.Exists(32)
```

## Methods
The `int` and `uint32` variations have the same API (except for the obvious difference that one deals with `int` and the other with `uint32`).

- `Set(int)`
- `Exsits(int) bool`
- `Remove(int) bool`

(It's hopefully obviously where a `uint32` is expected when dealing with the `uint32` variant)

