#ifndef ERROR_H
#define ERROR_H
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

typedef enum {
    ERROR_LEXER,
    ERROR_PARSER,
    ERROR_RUNTIME,
    ERROR_INTERNAL
} ErrorType;

typedef struct {
    char *message;
    int line;
    int column;
    const char *filename;
    const char *line_text;
    ErrorType type;
}  Error;

Error init_error(ErrorType type, const char *message, int line, int column, const char *filename, const char *line_text);
void print_error(Error *error);


#endif