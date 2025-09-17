package main

import "fmt"

func test() (int, string) {
	return 1, "zhangsan"
}

func main() {
	a, b := test()
	fmt.Println(a, b)

	a, _ = test()
	fmt.Println(a)

	_, b = test()
	fmt.Println(b)
}
