package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(greet("sanchu", 39))
	datas := make(map[string]string, 10)

	for i := 0; i < 10; i++ {
		datas[fmt.Sprint(i)] = fmt.Sprint("Sanchu_", i)
	}
	jsonString, err := json.Marshal(datas)
	fmt.Println(datas)
	fmt.Println(err)
	escape :=strconv.Quote(string(jsonString))
	fmt.Println(escape)

	fmt.Println("---------------------")
	s, _ := strconv.Unquote(escape)
	var msg map[string]string
	fmt.Println(s)
	err = json.Unmarshal(jsonString, &msg)
	fmt.Println(s)
	fmt.Println(msg)

}

func greet(name string, age int) string {
	return fmt.Sprint(name, " is ", age, " years old.")
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}
