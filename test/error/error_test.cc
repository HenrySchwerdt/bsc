#include <gtest/gtest.h>
#include <cstring>

extern "C" {
    #include "../../src/error/error.h"
}

TEST(ErrorTest, InitError) {
    const char * msg = "Unexpected character '#'";
    Error error = init_error(ERROR_LEXER,msg, 1, 1, "src/main.c", "#include <stdio.h>");
    ASSERT_EQ(error.type, ERROR_LEXER);
    ASSERT_STREQ(error.message, "Unexpected character '#'");
    ASSERT_EQ(error.line, 1);
    ASSERT_EQ(error.column, 1);
    ASSERT_STREQ(error.filename, "src/main.c");
    ASSERT_STREQ(error.line_text, "#include <stdio.h>");
}

TEST(ErrorTest, PrintLexingError) {
    const char * msg = "Unexpected character '#'";
    Error error = init_error(ERROR_LEXER, msg, 1, 1, "src/main.c", "#include <stdio.h>");
    print_error(&error);
}

TEST(ErrorTest, PrintParserError) {
    const char * msg = "Expected expression";
    Error error = init_error(ERROR_PARSER, msg, 1, 1, "src/main.c", "if");
    print_error(&error);
}


