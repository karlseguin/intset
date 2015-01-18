A specialized set for integers, ideal when:

- The number of elements is known ahead of time (or a good approximation)
- The number of elements doesn't change drastically over time
- The values are naturally random

As long as the number of elements within the set stays close to the originally specified size (I don't know the magic number, so let's say ±10%), and that they stay evenly distributed. the set will exhibit good read and write performance, as well as decent memory usage. When packed, read performance is roughly 2x better than a map[int]struct{} with memory usage less than 1/2.

