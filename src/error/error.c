#include "error.h"


Error init_error(ErrorType type, const char *message, int line, int column, const char *filename, const char *line_text) {
    Error error;
    error.type = type;
    error.message = message;
    error.line = line;
    error.column = column;
    error.filename = filename;
    error.line_text = line_text;
    return error;
}

/* 
Print errors with the following format:
Lexing error at src/main.c:1:1
1 | #include <stdio.h>
  | ^^^^^^^^^^^^^^^^^^
Unexpected character '#'
*/
void print_error(Error *error) {
    printf("\033[1;31m");
    switch (error->type) {
        case ERROR_LEXER:
            printf("Lexing error");
            break;
        case ERROR_PARSER:
            printf("Parsing error");
            break;
        case ERROR_COMPILER:
            printf("Compiler error");
            break;
        case ERROR_INTERNAL:
            printf("Internal error");
            break;
    }
    printf("\033[0m");
    printf(" at %s:%d:%d\n", error->filename, error->line, error->column);
    // add line numbers with light blue color
    printf("\033[1;36m");
    printf("%d | %s\n", error->line, error->line_text);
    printf("  | ");
    for (int i = 0; i < error->column - 1; i++) {
        printf(" ");
    }
    printf("\033[0;31m");
    for (int i = 0; i < strlen(error->line_text); i++) {
        printf("^");
    }
    printf("\n");
    printf("\033[0m");
    printf("%s\n", error->message);
}

