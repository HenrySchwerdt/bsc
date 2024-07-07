#ifndef ERROR_H
#define ERROR_H
#include <stdio.h>
#include <string.h>

typedef enum {
    ERROR_LEXER,
    ERROR_PARSER,
    ERROR_COMPILER,
    ERROR_INTERNAL
} ErrorType;

typedef struct {
    const char *message;
    int line;
    int column;
    const char *filename;
    const char *line_text;
    ErrorType type;
}  Error;

Error init_error(ErrorType type, const char *message, int line, int column, const char *filename, const char *line_text);
void print_error(Error *error);


#endif