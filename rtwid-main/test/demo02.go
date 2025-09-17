package main

import "fmt"

var (
	people   string
	age      int
	location string
)

func main() {

	people = "Tom"
	age = 18
	location = "China"

	fmt.Println("Name:", people, "Age:", age, "Location:", location)
}
