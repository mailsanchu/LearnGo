package main

import "fmt"

func main() {
	safe, speed := true, 50
	fmt.Println(safe, speed)

	var safeS bool
	safeS, demoX := false, "Demo"
	fmt.Println(safeS, demoX)

	var (
		name   string
		age    int
		famous bool
	)
	name = "Sanchu"
	age = 39
	famous = true
	//fmt.Println(name, age, famous)

	name = "Jubina"
	age = 38
	famous = true
	fmt.Println(name, age, famous)

	var Test1 string
	Test1 = "Sanchu varkey"

	var Test2 *string
	Test2 = &Test1
	Test1 = "Sanchu"
	Test1 = "Demo"
	Test1 = "Sanchu Varkey Key"
	fmt.Println(*Test2)
	fmt.Println(Test1)

	fmt.Println(Reverse("The quick brown $ jumped over the lazy 犬"))
	fmt.Println(Reverse1("The quick brown $ jumped over the lazy 犬"))
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Reverse1(s string) (result string) {
	for _,v := range s {
		result = string(v) + result
	}
	return
}
