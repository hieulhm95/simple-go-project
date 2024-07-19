// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	fmt.Println(variadicSum(1, 2, 3, 4, 45))
}

func variadicSum(nums ...int) (sum int) {
	for _, v := range nums {
		sum += v
	}
	return
}
