package main

import "fmt"

func main() {
	const (
		a = iota
		b = iota
		c
		d = "haha"
		e
		f = 100
		g
		h = iota
		i
	)

	fmt.Println(a, b, c, d, e, f, g, h, i)
}
