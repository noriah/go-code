lst = [1, "Fizz", "Buzz", "FizzBuzz"]
for (; lst[0] <= 100; lst[0]++) {
  divis3 = 1 >> (lst[0] % 3)
  divis5 = (1 >> (lst[0] % 5)) << 1
  console.log(lst[divis3 | divis5])
}
