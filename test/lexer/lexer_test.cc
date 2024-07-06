#include <gtest/gtest.h>
#include <cstring>

extern "C" {
    #include "../../src/lexer/lexer.h"
}

void compare_token(const Token& token, TokenType expected_type, const char* expected_value, int expected_line, int expected_column) {
    ASSERT_EQ(token.type, expected_type);
    ASSERT_STREQ(token.value, expected_value);
    ASSERT_EQ(token.line, expected_line);
    ASSERT_EQ(token.column, expected_column);
}

TEST(LexerTest, HandlesSingleCharacterTokens) {
    const char* source = "(){}[],;:.?+-*/%&|~!";
    Lexer lexer;
    init_lexer(&lexer, source);
    init_transition_table();

    compare_token(next_token(&lexer), TK_LPAREN, "(", 1, 1);
    compare_token(next_token(&lexer), TK_RPAREN, ")", 1, 2);
    compare_token(next_token(&lexer), TK_LBRACE, "{", 1, 3);
    compare_token(next_token(&lexer), TK_RBRACE, "}", 1, 4);
    compare_token(next_token(&lexer), TK_LBRACKET, "[", 1, 5);
    compare_token(next_token(&lexer), TK_RBRACKET, "]", 1, 6);
    compare_token(next_token(&lexer), TK_COMMA, ",", 1, 7);
    compare_token(next_token(&lexer), TK_SEMICOLON, ";", 1, 8);
    compare_token(next_token(&lexer), TK_COLON, ":", 1, 9);
    compare_token(next_token(&lexer), TK_DOT, ".", 1, 10);
    compare_token(next_token(&lexer), TK_QUESTION, "?", 1, 11);
    compare_token(next_token(&lexer), TK_PLUS, "+", 1, 12);
    compare_token(next_token(&lexer), TK_MINUS, "-", 1, 13);
    compare_token(next_token(&lexer), TK_STAR, "*", 1, 14);
    compare_token(next_token(&lexer), TK_SLASH, "/", 1, 15);
    compare_token(next_token(&lexer), TK_PERCENT, "%", 1, 16);
    compare_token(next_token(&lexer), TK_AMPERSAND, "&", 1, 17);
    compare_token(next_token(&lexer), TK_PIPE, "|", 1, 18);
    compare_token(next_token(&lexer), TK_TILDE, "~", 1, 19);
    compare_token(next_token(&lexer), TK_BANG, "!", 1, 20);
    compare_token(next_token(&lexer), TK_EOF, "", 1, 20);
}

TEST(LexerTest, HandleKeywords) {
    const char* source = "if else while for return match let const class interface abstract static";
    Lexer lexer;
    init_lexer(&lexer, source);
    init_transition_table();

    compare_token(next_token(&lexer), TK_IF, "if", 1, 1);
    compare_token(next_token(&lexer), TK_ELSE, "else", 1, 4);
    compare_token(next_token(&lexer), TK_WHILE, "while", 1, 9);
    compare_token(next_token(&lexer), TK_FOR, "for", 1, 15);
    compare_token(next_token(&lexer), TK_RETURN, "return", 1, 19);
    compare_token(next_token(&lexer), TK_MATCH, "match", 1, 26);
    compare_token(next_token(&lexer), TK_LET, "let", 1, 32);
    compare_token(next_token(&lexer), TK_CONST, "const", 1, 36);
    compare_token(next_token(&lexer), TK_CLASS, "class", 1, 42);
    compare_token(next_token(&lexer), TK_INTERFACE, "interface", 1, 48);
    compare_token(next_token(&lexer), TK_ABSTRACT, "abstract", 1, 58);
    compare_token(next_token(&lexer), TK_STATIC, "static", 1, 67);
    compare_token(next_token(&lexer), TK_EOF, "", 1, 72);
}

TEST(LexerTest, HandlesIdentifiers) {
    const char* source = "variable_name anotherVariable yetAnotherVar";
    Lexer lexer;
    init_lexer(&lexer, source);
    init_transition_table();

    compare_token(next_token(&lexer), TK_IDENTIFIER, "variable_name", 1, 1);
    compare_token(next_token(&lexer), TK_IDENTIFIER, "anotherVariable", 1, 15);
    compare_token(next_token(&lexer), TK_IDENTIFIER, "yetAnotherVar", 1, 31);
    compare_token(next_token(&lexer), TK_EOF, "", 1, 43);
}

TEST(LexerTest, HandlesNumbers) {
    const char* source = "123 4567 89";
    Lexer lexer;
    init_lexer(&lexer, source);
    init_transition_table();

    compare_token(next_token(&lexer), TK_NUMBER, "123", 1, 1);
    compare_token(next_token(&lexer), TK_NUMBER, "4567", 1, 5);
    compare_token(next_token(&lexer), TK_NUMBER, "89", 1, 10);
    compare_token(next_token(&lexer), TK_EOF, "", 1, 11);
}

TEST(LexerTest, HandlesStrings) {
    const char* source = "\"hello\" \"world\"";
    Lexer lexer;
    init_lexer(&lexer, source);
    init_transition_table();

    compare_token(next_token(&lexer), TK_STRING, "hello", 1, 2);
    compare_token(next_token(&lexer), TK_STRING, "world", 1, 10);
    // compare_token(next_token(&lexer), TK_EOF, "", 1, 16);
}

TEST(LexerTest, HandlesCharacters) {
    const char* source = "'a' 'b'";
    Lexer lexer;
    init_lexer(&lexer, source);
    init_transition_table();

    compare_token(next_token(&lexer), TK_CHAR, "a", 1, 2);
    compare_token(next_token(&lexer), TK_CHAR, "b", 1, 6);
    // compare_token(next_token(&lexer), TK_EOF, "", 1, 9);
}