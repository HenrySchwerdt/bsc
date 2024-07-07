#ifndef RESULT_H
#define RESULT_H

#include <stdlib.h>


typedef void (*ResolveFunc)(void *value);
typedef void (*RejectFunc)(void *error);

typedef struct Result {
    void *value;
    void *error;
} Result;

Result ok(void *value);
Result error(void *error);
void resolve(Result *result, ResolveFunc resolve, RejectFunc reject);


#endif // RESULT_H
