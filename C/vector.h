#ifndef VECTOR_H
#define VECTOR_H

typedef struct Vector {
  int cap;
  int length;
  long long* values;
} Vector;

typedef struct Index {
  int index;
  int exists;
} Index;

Vector* vector_new();
Index vector_index(Vector* vector, long long value);
void vector_insert(Vector* vector, int index, long long value);
void vector_free(Vector*);

#endif
