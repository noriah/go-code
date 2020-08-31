function fizzBuzzArray() {
  var arr = [1, "Fizz", "Buzz", "FizzBuzz"]
  for (; arr[0] <= 100; arr[0]++) {
    div3 = 1 >> (arr[0] % 3)
    div5 = (1 >> (arr[0] % 5)) << 1
    console.log(arr[div3 | div5])
  }
}

function fizzBuzzArrayAlt() {
  var div3, div5
  var arr = [1, "Fizz", "Buzz", "FizzBuzz"]
  for (; arr[0] <= 100; arr[0]++) {
    div3 = 1 >> (arr[0] % 3)
    div5 = (1 >> (arr[0] % 5)) << 1
    console.log(arr[div3 | div5])
  }
}

fizzBuzzArray()
