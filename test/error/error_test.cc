#include <gtest/gtest.h>
#include <cstring>

extern "C" {
    #include "../../src/error/error.h"
}

TEST(ErrorTest, InitError) {
    Error error = init_error(ERROR_LEXER, "Unexpected character '#'", 1, 1, "src/main.c", "#include <stdio.h>");
    ASSERT_EQ(error.type, ERROR_LEXER);
    ASSERT_STREQ(error.message, "Unexpected character '#'");
    ASSERT_EQ(error.line, 1);
    ASSERT_EQ(error.column, 1);
    ASSERT_STREQ(error.filename, "src/main.c");
    ASSERT_STREQ(error.line_text, "#include <stdio.h>");
    print_error(&error);
}