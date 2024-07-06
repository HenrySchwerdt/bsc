#ifndef LEXER_H
#define LEXER_H
#include <string.h>
#include <stdlib.h>
typedef enum
{
    TK_EOF,
    // SINGLE CHAR TOKENS
    TK_LPAREN,
    TK_RPAREN,
    TK_LBRACE,
    TK_RBRACE,
    TK_LBRACKET,
    TK_RBRACKET,
    TK_COMMA,
    TK_DOT,
    TK_SEMICOLON,
    TK_COLON,
    TK_QUESTION,
    TK_PLUS,
    TK_MINUS,
    TK_STAR,
    TK_SLASH,
    TK_PERCENT,
    TK_CARET,
    TK_AMPERSAND,
    TK_PIPE,
    TK_TILDE,
    TK_BANG,
    TK_LESS,
    TK_GREATER,
    TK_EQUAL,

    // DOUBLE CHAR TOKENS
    TK_PLUS_PLUS,
    TK_MINUS_MINUS,
    TK_PLUS_EQUAL,
    TK_MINUS_EQUAL,
    TK_STAR_EQUAL,
    TK_SLASH_EQUAL,
    TK_PERCENT_EQUAL,
    TK_CARET_EQUAL,
    TK_AMPERSAND_EQUAL,
    TK_TILDE_EQUAL,
    TK_BANG_EQUAL,
    TK_LESS_EQUAL,
    TK_GREATER_EQUAL,
    TK_EQUAL_EQUAL,
    TK_LESS_LESS,
    TK_GREATER_GREATER,
    TK_LESS_LESS_EQUAL,
    TK_GREATER_GREATER_EQUAL,
    TK_AMPERSAND_AMPERSAND,
    TK_PIPE_PIPE,
    TK_ARROW,

    // Keywords
    TK_IF,
    TK_ELSE,
    TK_WHILE,
    TK_FOR,
    TK_RETURN,
    TK_MATCH,
    TK_LET,
    TK_CONST,
    TK_CLASS,
    TK_INTERFACE,
    TK_ABSTRACT,
    TK_STATIC,

    // Values
    TK_IDENTIFIER,
    TK_NUMBER,
    TK_STRING,
    TK_CHAR,

} TokenType;

typedef struct {
    TokenType type;
    char *value;
    int line;
    int column;
} Token;


typedef enum
{
    STATE_START,
    STATE_IDENTIFIER,
    STATE_NUMBER,
    STATE_STRING,
    STATE_CHAR,
    STATE_OPERATOR,
    STATE_FINAL,
    STATE_STRING_END,
    STATE_CHAR_END,
    STATE_ERROR
} LexerState;

typedef struct {
    const char *input;
    int position;
    int line;
    int column;
    int length;
    LexerState state;
} Lexer;

Token next_token(Lexer *lexer);
void init_lexer(Lexer *lexer, const char *input);
void init_transition_table();

#endif