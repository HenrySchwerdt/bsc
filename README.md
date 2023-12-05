# bsc

# BNF
```
<program> ::= <statement>*

<statement> ::= <declaration-statement>
            | <ternary-statement>
            | <assignment-statement>
            | <call-statement>
            | <while-statement>
            | <for-statement>
            | <if-statement>
            | <block-statement>
            | <fn-declaration-statement>
            | <import-statement>
            | <break-statement>
            | <continue-statement>
            | <return-statement>

# Statements
<import-statement> ::= 'import' '{' IDENTIFIER (, IDENTIFIER)* '}' 'from' '\''PATH'\''
<fn-declaration-statement> ::= 'fn' IDENTIFIER '(' <param>? (',' <param>)* ')' ':' <type> <block-statement>
<if-statement> ::= 'if' '(' <expression> ')' <block-statement> ('else' <block-statement>)?
<while-statement> ::= 'while' '(' <expression> ')' <block-statement>
<for-statement> ::= 'for' '(' (<var-declaration> | <assignment-expression>)? ';' <expression> ';' <assignment-expression> ')' <block-statement>
<block-statement> ::= '{' <statement>* '}'
<call-statement> ::= <call-expression>';'
<ternary-statement> ::= <ternary-expression>';'
<assignment-statement> ::= <assignment-expression>';'
<declaration-statement> ::= <var-declaration-expression>';'
<break-statement> ::= 'break' ';'
<continue-statement> ::= 'continue' ';'
<return-statement> ::= 'return' <expression> ';'


# Expressions
<assignment-expression> ::= IDENTIFIER '=' <expression>

<var-declaration-expression> ::= <specifier> IDENTIFIER (: <type>)? '=' <expression> | <specifier> IDENTIFIER : <type>
<expression> ::= <logical-or-expression>

<logical-or-expression> ::= <logical-and-expression> ('||' <logical-and-expression>)*
<logical-and-expression> ::= <equality-expression> ('&&' <equality-expression>)*
<equality-expression> ::= <relational-expression> (('==' | '!=') <relational-expression>)*
<relational-expression> ::= <additive-expression> (('<' | '>' | '<=' | '>=') <additive-expression>)*
<additive-expression> ::= <multiplicative-expression> (('+' | '-') <multiplicative-expression>)*
<multiplicative-expression> ::= <unary-expression> (('*' | '/' | '%') <unary-expression>)*
<unary-expression> ::= ('+' | '-' | '!')* <primary-expression>

<primary-expression> ::= IDENTIFIER
                | LITERAL
                | '(' <expression> ')'
                | <call-expression>
<call-expression> ::= IDENTIFIER '(' <expression>? (,<expression>)* ')'

# Util
<specifier> ::= 'let' | 'const'
<type> ::= 'int8' | 'int16' | 'int32' | 'int64' | 'uint8' | 'uint16' | 'uint32' | 'uint64' | 'bool' | 'void' | 'float32' | 'float64' | <complex-type>
<complex-type> ::= IDENTIFIER ('<' <complex-type> '>')? || '(' <type>? (',' <type>)* ')' '=>' <type>
<param> ::= IDENTIFIER ':' <type>


Program = Statement* .
Statement = VarDeclarationStatement 
            | AssignmentStatement 
            | FnDeclarationStatement 
            | ImportStatement 
            | BlockStatement 
            | IfStatement 
            | WhileStatement 
            | ForStatement 
            | BreakStatement 
            | ReturnStatement 
            | ContinueStatement 
            | ExitStatement 
            | CallStatement .
VarDeclarationStatement = VarDeclarationExpression ";" .
VarDeclarationExpression = ("let" | "const") <ident> ":" Type? ("=" Expression)? ("=" ArrayInitializer)? .
Type = ("int8" | "int16" | "int32" | "int64" | "uint8" | "uint16" | "uint32" | "uint64" | "bool" | "void" | "float32" | "float64") ("[" "]")* .
Expression = LogicalOrExpression .
LogicalOrExpression = LogicalAndExpression ("||" LogicalOrExpression)? .
LogicalAndExpression = EqualityExpression ("&&" LogicalAndExpression)? .
EqualityExpression = RelationalExpression ("==" | "!=")? EqualityExpression? .
RelationalExpression = AdditiveExpression ("<" | ">" | "<=" | ">=")? RelationalExpression? .
AdditiveExpression = MultiplicativeExpression ("+" | "-")? AdditiveExpression? .
MultiplicativeExpression = UnaryExpression ("*" | "/" | "%")? MultiplicativeExpression? .
UnaryExpression = PrimaryExpression .
PrimaryExpression = CallExpression | Value | ArrayLookup | ("(" Expression ")") | <ident> .
CallExpression = <ident> "(" Expression? ("," Expression)* ")" .
Value = Float | Int | String | Bool .
Float = <float> .
Int = <int> .
String = <string> .
Bool = ("true" | "false") .
ArrayLookup = <ident> "[" Expression "]" .
ArrayInitializer = "[" Value? ("," Value)* "]" .
AssignmentStatement = AssignmentExpression ";" .
AssignmentExpression = <ident> ("=" | "+=" | "-=" | "*=" | "/=" | "|=") Expression? ArrayInitializer? .
FnDeclarationStatement = "fn" <ident> "(" Param? ("," Param)* ")" ":" Type BlockStatement .
Param = <ident> ":" Type .
BlockStatement = "{" Statement* "}" .
ImportStatement = "import" "{" <ident> ("," <ident>)* "}" "from" <string> .
IfStatement = "if" "(" Expression ")" BlockStatement ("else" BlockStatement)? .
WhileStatement = "while" "(" Expression ")" BlockStatement .
ForStatement = "for" "(" VarDeclarationExpression? ";" Expression? ";" AssignmentExpression? ")" BlockStatement .
BreakStatement = "break" ";" .
ReturnStatement = "return" Expression ";" .
ContinueStatement = "continue" ";" .
ExitStatement = "exit" "(" Expression ")" ";" .
CallStatement = CallExpression ";" .

```


## Examples
### Loops
```bs
// Example for For-Loop in BlockScript
var sum = 0;
for (var i = 0; i<=5; i += 1) {
    sum = sum + 1;
    for (var j = 0; j <= 5; j += 1) {
        sum = sum + j;
    }   
}
exit(sum);
```
```bs
// Example for While-Loop in BlockScript
var x = 0;
var sum = 0;
while (x <= 5) {
    var xI = 0;
    while(xI <= 5) {
        xI += 1;
        sum += xI;
    }
    x += 1;
    sum += x;
}
exit(sum);
```

### Rekursion
Currently you can only create int functions that return an integer. Soon there will be other <type>s and build in std.
```bs
// Example for Fib in BlockScript
fn fib(n) Int {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
}
exit(fib(6));
```
```bs
// Example for Factorial in BlockScript
fn factorial(n) Int {
    if (n == 0) return 1;
    return n * factorial(n - 1);
}
exit(factorial(4));
```
```bs
// Example for GCD in BlockScript
fn gcd(a, b) Int {
    if (b == 0) return a;
    return gcd(b, a % b);
}
exit(gcd(56, 98));
```