#include "lexer.h"
#define NUM_STATES 10
#define NUM_CHARS 128

static const char *keywords[] = {
    "if", "else", "while", "for", "return", "match", "let", "const",
    "class", "interface", "abstract", "static", NULL};

static TokenType keyword_tokens[] = {
    TK_IF, TK_ELSE, TK_WHILE, TK_FOR, TK_RETURN, TK_MATCH, TK_LET, TK_CONST,
    TK_CLASS, TK_INTERFACE, TK_ABSTRACT, TK_STATIC};

LexerState transition[NUM_STATES][NUM_CHARS];

void init_transition_table()
{
    memset(transition, STATE_ERROR, sizeof(transition));
    // Sets states for all identifiers
    for (char c = 'a'; c <= 'z'; c++)
    {
        transition[STATE_START][(int)c] = STATE_IDENTIFIER;
        transition[STATE_IDENTIFIER][(int)c] = STATE_IDENTIFIER;
    }
    for (char c = 'A'; c <= 'Z'; c++)
    {
        transition[STATE_START][(int)c] = STATE_IDENTIFIER;
        transition[STATE_IDENTIFIER][(int)c] = STATE_IDENTIFIER;
    }
    transition[STATE_START]['_'] = STATE_IDENTIFIER;
    transition[STATE_IDENTIFIER]['_'] = STATE_IDENTIFIER;

    for (char c = '0'; c <= '9'; c++)
    {
        transition[STATE_START][(int)c] = STATE_NUMBER;
        transition[STATE_NUMBER][(int)c] = STATE_NUMBER;
    }

    // Add transitions for single character tokens
    transition[STATE_START]['('] = STATE_FINAL;
    transition[STATE_START][')'] = STATE_FINAL;
    transition[STATE_START]['{'] = STATE_FINAL;
    transition[STATE_START]['}'] = STATE_FINAL;
    transition[STATE_START]['['] = STATE_FINAL;
    transition[STATE_START][']'] = STATE_FINAL;
    transition[STATE_START][','] = STATE_FINAL;
    transition[STATE_START]['.'] = STATE_FINAL;
    transition[STATE_START][';'] = STATE_FINAL;
    transition[STATE_START][':'] = STATE_FINAL;
    transition[STATE_START]['?'] = STATE_FINAL;
    transition[STATE_START]['+'] = STATE_OPERATOR;
    transition[STATE_START]['-'] = STATE_OPERATOR;
    transition[STATE_START]['*'] = STATE_FINAL;
    transition[STATE_START]['/'] = STATE_FINAL;
    transition[STATE_START]['%'] = STATE_FINAL;
    transition[STATE_START]['^'] = STATE_FINAL;
    transition[STATE_START]['&'] = STATE_OPERATOR;
    transition[STATE_START]['|'] = STATE_OPERATOR;
    transition[STATE_START]['~'] = STATE_FINAL;
    transition[STATE_START]['!'] = STATE_OPERATOR;
    transition[STATE_START]['<'] = STATE_OPERATOR;
    transition[STATE_START]['>'] = STATE_OPERATOR;
    transition[STATE_START]['='] = STATE_OPERATOR;

    // Add transitions for strings
    transition[STATE_START]['"'] = STATE_STRING;
    for (char c = 32; c <= 126; c++)
    {
        transition[STATE_STRING][(int)c] = STATE_STRING;
    }
    transition[STATE_STRING]['"'] = STATE_STRING_END;

    // Add transitions for characters
    transition[STATE_START]['\''] = STATE_CHAR;
    for (char c = 32; c <= 126; c++)
    {
        transition[STATE_CHAR][(int)c] = STATE_CHAR;
    }
    transition[STATE_CHAR]['\''] = STATE_CHAR_END;
}

static Token make_token(Lexer *lexer, TokenType type, const char *start, int length)
{
    Token token;
    token.type = type;
    token.value = strndup(start, length);
    token.line = lexer->line;
    token.column = lexer->column - length;
    return token;
}

static TokenType is_keyword(const char *str, int length)
{
    for (int i = 0; keywords[i] != NULL; i++)
    {
        if (strncmp(str, keywords[i], length) == 0)
        {
            return keyword_tokens[i];
        }
    }
    return TK_IDENTIFIER;
}

static int is_whitespace(char c)
{
    return c == ' ' || c == '\t' || c == '\n';
}

static int is_final_state(LexerState state)
{
    return state == STATE_FINAL || state == STATE_STRING_END || state == STATE_CHAR_END;
}

Token next_token(Lexer *lexer)
{
    while (is_whitespace(lexer->input[lexer->position]))
    {
        if (lexer->input[lexer->position] == '\n')
        {
            lexer->line++;
            lexer->column = 0;
        }
        else
        {
            lexer->column++;
        }
        lexer->position++;
    }
    lexer->state = STATE_START;
    const char *start = &lexer->input[lexer->position];
    char current_char;
    lexer->length = 0;
    while (!is_final_state(lexer->state))
    {
        current_char = lexer->input[lexer->position];
        if (current_char == '\n')
        {
            lexer->line++;
            lexer->column = 0;
            continue;
        }
        if (current_char == '\0')
        {
            break;
        }
        if (is_whitespace(current_char) && lexer->state != STATE_STRING && lexer->state != STATE_CHAR)
        {
            break;
        }
        LexerState next_state = transition[lexer->state][(int)current_char];
        if (lexer->state != STATE_START && lexer->state != next_state && next_state != STATE_FINAL && next_state != STATE_OPERATOR)
        {
            break;
        }
        if (lexer->state == STATE_NUMBER && next_state != STATE_NUMBER)
        {
            break;
        }
        if (lexer->state == STATE_OPERATOR && next_state != STATE_OPERATOR)
        {
            break;
        }

        lexer->state = next_state;
        lexer->length++;
        lexer->position++;
        lexer->column++;
    }
    switch (lexer->state)
    {
    case STATE_IDENTIFIER:
        return make_token(lexer, is_keyword(start, lexer->position - (start - lexer->input)), start, lexer->position - (start - lexer->input));
    case STATE_NUMBER:
        return make_token(lexer, TK_NUMBER, start, lexer->position - (start - lexer->input));
    case STATE_STRING_END:
        return make_token(lexer, TK_STRING, start + 1, lexer->position - (start - lexer->input) - 2);
    case STATE_CHAR_END:
        return make_token(lexer, TK_CHAR, start + 1, lexer->position - (start - lexer->input) - 2);
    case STATE_OPERATOR:
        if (lexer->length > 1)
        {
            if (lexer->input[lexer->position - 1] == '+' && lexer->input[lexer->position] == '+')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_PLUS_PLUS, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '-' && lexer->input[lexer->position] == '-')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_MINUS_MINUS, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '&' && lexer->input[lexer->position] == '&')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_AMPERSAND_AMPERSAND, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '|' && lexer->input[lexer->position] == '|')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_PIPE_PIPE, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '=' && lexer->input[lexer->position] == '>')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_ARROW, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '<' && lexer->input[lexer->position] == '<')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_LESS_LESS, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '>' && lexer->input[lexer->position] == '>')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_GREATER_GREATER, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '<' && lexer->input[lexer->position] == '=')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_LESS_EQUAL, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '>' && lexer->input[lexer->position] == '=')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_GREATER_EQUAL, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '>' && lexer->input[lexer->position] == '>')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_GREATER_GREATER, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '<' && lexer->input[lexer->position] == '<')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_LESS_LESS, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '=' && lexer->input[lexer->position] == '=')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_EQUAL_EQUAL, start, 2);
            }
            else if (lexer->input[lexer->position - 1] == '!' && lexer->input[lexer->position] == '=')
            {
                lexer->position++;
                lexer->column++;
                return make_token(lexer, TK_BANG_EQUAL, start, 2);
            }
        }
        else
        {
            switch (*start)
            {
            case '+':
                return make_token(lexer, TK_PLUS, start, 1);
            case '-':
                return make_token(lexer, TK_MINUS, start, 1);
            case '&':
                return make_token(lexer, TK_AMPERSAND, start, 1);
            case '|':
                return make_token(lexer, TK_PIPE, start, 1);
            case '=':
                return make_token(lexer, TK_EQUAL, start, 1);
            case '<':
                return make_token(lexer, TK_LESS, start, 1);
            case '>':
                return make_token(lexer, TK_GREATER, start, 1);
            case '!':
                return make_token(lexer, TK_BANG, start, 1);
            default:
                return make_token(lexer, TK_EOF, start, 1);
            }
        }
    case STATE_FINAL:
        switch (*start)
        {
        case '(':
            return make_token(lexer, TK_LPAREN, start, 1);
        case ')':
            return make_token(lexer, TK_RPAREN, start, 1);
        case '{':
            return make_token(lexer, TK_LBRACE, start, 1);
        case '}':
            return make_token(lexer, TK_RBRACE, start, 1);
        case '[':
            return make_token(lexer, TK_LBRACKET, start, 1);
        case ']':
            return make_token(lexer, TK_RBRACKET, start, 1);
        case ',':
            return make_token(lexer, TK_COMMA, start, 1);
        case '.':
            return make_token(lexer, TK_DOT, start, 1);
        case ';':
            return make_token(lexer, TK_SEMICOLON, start, 1);
        case ':':
            return make_token(lexer, TK_COLON, start, 1);
        case '?':
            return make_token(lexer, TK_QUESTION, start, 1);
        case '*':
            return make_token(lexer, TK_STAR, start, 1);
        case '/':
            return make_token(lexer, TK_SLASH, start, 1);
        case '%':
            return make_token(lexer, TK_PERCENT, start, 1);
        case '^':
            return make_token(lexer, TK_CARET, start, 1);
        case '~':
            return make_token(lexer, TK_TILDE, start, 1);
        case '!':
            return make_token(lexer, TK_BANG, start, 1);
        case '<':
            return make_token(lexer, TK_LESS, start, 1);
        case '>':
            return make_token(lexer, TK_GREATER, start, 1);
        case '=':
            return make_token(lexer, TK_EQUAL, start, 1);
        default:
            return make_token(lexer, TK_EOF, start, 1);
        }
    default:
        return make_token(lexer, TK_EOF, start, 1);
    }
}

void init_lexer(Lexer *lexer, const char *input)
{
    lexer->input = input;
    lexer->position = 0;
    lexer->length = 0;
    lexer->line = 1;
    lexer->column = 1;
    lexer->state = STATE_START;
}