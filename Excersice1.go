package main

import "fmt"

func main() {
	h := half(78, 89, 45, 56, 14 ,1500 ,7777 , 9555.0)
	fmt.Println(h)
	var sliceEx []string
	sliceEx = append(sliceEx, "Sanchu")
	fmt.Println(&sliceEx)
	fmt.Printf("address of slice %p add of Arr %p \n", &sliceEx, &sliceEx)
	sliceEx = append(sliceEx, "Sanchu")
	fmt.Println(&sliceEx)
	fmt.Printf("address of slice %p add of Arr %p \n", &sliceEx, &sliceEx)

}

func half(n ... int) int {
	temp := 0
	for _, val := range n {
		if val > temp {
			temp = val
		}
	}
	return temp
}
