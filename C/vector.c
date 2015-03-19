#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "vector.h"

Vector* vector_new() {
	Vector* vector = (Vector*) malloc(sizeof(Vector));
	if (!vector) {
		return NULL;
	}
	vector->cap = 8;
	vector->length = 0;
	vector->values = (long long*)malloc(sizeof(long long) * vector->cap);
	return vector;
}

Index vector_index(Vector* vector, long long value) {
	int l = vector->length;
	Index index = {l, 0};
	for (int i = 0; i < l; i++) {
		int v = vector->values[i];
		if (v < value) continue;
		index.index = i;
		index.exists = v == value ? 1 : 0;
		return index;
	}
	return index;
}

void vector_insert(Vector *vector, int index, long long value) {
	int length = vector->length;
	if (length == vector->cap) {
		vector->cap += 8;
		//todo: handle null
		vector->values = (long long*) realloc(vector->values, sizeof(long long) * vector->cap);
	}

	if (index != length) {
		for (int i = length - 1; i >= index; i--) {
			vector->values[i+1] = vector->values[i];
		}
	}
	vector->values[index] = value;
	vector->length++;
}

void vector_free(Vector* vector) {
	free(vector->values);
}
