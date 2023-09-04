# bsc

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
Currently you can only create int functions that return an integer. Soon there will be other types and build in std.
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