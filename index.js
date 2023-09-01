var x = 0;
var sum = 0;
while (x <= 1) {
  var innerX = 0;
  while (innerX <= 1) {
    innerX = innerX + 1;
    sum = sum + innerX;
  }
  x = x + 1;
  sum = sum + x;
}
console.log(sum);
