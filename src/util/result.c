#include "result.h"

Result ok(void *value) {
    Result result;
    result.value = value;
    result.error = NULL;
    return result;
}

Result error(void *error) {
    Result result;
    result.value = NULL;
    result.error = error;
    return result;
}

void resolve(Result *result, ResolveFunc resolve, RejectFunc reject) {
    if (result->value == NULL) {
        reject(result->error);
    } else {
        resolve(result->value);
    }
}
