var prev = 0;
var next = 1;
var nFib = 10;
var it = 1;
while(it < nFib) {
    var sum = prev + next;
    prev = next;
    next = sum;
    it = it + 1;
}
console.log(prev);