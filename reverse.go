package main

import (
	"fmt"
	"math"
)

func main() {
	reversedInteger := reverse(152)
	fmt.Println(reversedInteger)

	reversedInteger = reverse(-152)
	fmt.Println(reversedInteger)

	reversedInteger = reverse(290)
	fmt.Println(reversedInteger)
}

func reverse(x int) int {
	sign := "positive"
	if x >= 0 {
		sign = "positive"
	} else {
		sign = "negative"
	}

	x = int(math.Abs(float64(x)))

	var reversedDigit int

	for x > 0 {
		lastDigit := x % 10
		reversedDigit = reversedDigit*10 + lastDigit

		x = x / 10
	}

	if sign == "negative" {
		reversedDigit = reversedDigit * -1
	}
	return reversedDigit
}
