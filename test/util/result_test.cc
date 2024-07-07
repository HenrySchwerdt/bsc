#include <gtest/gtest.h>
#include <cstring>

extern "C" {
    #include "../../src/util/result.h"
}

TEST(ResultTest, Resolve) {
    Result result = ok((void *)"Hello, World!");
    resolve(&result, [](void *value) {
        ASSERT_STREQ((char *)value, "Hello, World!");
    }, [](void* error) {
        FAIL();
    });
}

TEST(ResultTest, Reject) {
    Result result = error((void *)"An error occurred");
    resolve(&result, [](void *value) {
        FAIL();
    }, [](void* error) {
        ASSERT_STREQ((char *)error, "An error occurred");
    });
}
