#include <stdio.h>
#include <stdlib.h>
#include "intset.h"
#include "vector.h"

IntSet* intset_new(long long size){
  if (size < BUCKET_SIZE) {
    size = BUCKET_SIZE * 2;
  }
  int bucketCount = upTwo(size) / BUCKET_SIZE;

  IntSet *set = (IntSet*) malloc(sizeof(IntSet));
  if (!set) goto err;

  set->length = 0;
  set->mask = bucketCount - 1;
  set->buckets = (Vector**) malloc(sizeof(Vector*) * bucketCount);
  if (!set->buckets) goto err;
  for (int i = 0; i < bucketCount; i++) {
    //todo: handle NULL
    set->buckets[i] = vector_new();
  }
  return set;

  err:
    if (set) intset_free(set);
    return NULL;
}

void intset_set(IntSet* set, long long value) {
  Vector* bucket = set->buckets[value & set->mask];
  Index index = vector_index(bucket, value);
  if (index.exists) return;
  vector_insert(bucket, index.index, value);
  set->length++;
}

int intset_exists(IntSet* set, long long value) {
  Vector* bucket = set->buckets[value & set->mask];
  return vector_index(bucket, value).exists;
}

IntSet* intset_intersect(IntSet* set1, IntSet* set2) {
  int count = 0;
  int values[set1->length];
  int bucketCount = set1->mask + 1;
  for (int i = 0; i < bucketCount; i++) {
    Vector *bucket = set1->buckets[i];
    for (int j = 0; j < bucket->length; j++) {
      int value = bucket->values[j];
      if (intset_exists(set2, value)) {
        values[count] = value;
        count++;
      }
    }
  }

  IntSet* new = intset_new(count);
  for (int i = 0; i < count; i++) {
    intset_set(new, values[i]);
  }
  return new;
}

void intset_free(IntSet* set) {
  int bucketCount = set->mask + 1;
  for (int i = 0; i < bucketCount; i++) {
    vector_free(set->buckets[i]);
  }
  free(set->buckets);
  free(set);
}


// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
int upTwo(int v) {
	v--;
	v |= v >> 1;
	v |= v >> 2;
	v |= v >> 4;
	v |= v >> 8;
	v |= v >> 16;
	return ++v;
}
