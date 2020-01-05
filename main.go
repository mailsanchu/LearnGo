package main

import "fmt"

const (
	OK = "Jubina"
)

func Hello() {
	//const OK = true
	var HelloV = "Sanchu"
	DEMO := &HelloV
	*DEMO = "Java"
	fmt.Println(OK, HelloV, "Demo", *DEMO)
}

func main() {
	Hello()
}
