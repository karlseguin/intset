#ifndef INTSET_H
#define INTSET_H

#include "vector.h"

#define BUCKET_SIZE 32

typedef struct IntSet{
  int mask;
  int length;
  Vector** buckets;
} IntSet;

IntSet* intset_new(long long size);
void intset_set(IntSet* set, long long value);
int intset_exists(IntSet* set, long long value);
IntSet* intset_intersect(IntSet* set1, IntSet* set2);
void intset_free(IntSet* set);

int upTwo(int value);

#endif
