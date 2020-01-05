package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// range specification, note that min <= max
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

func stringToBin(s string) (binString string) {
	for _, c := range s {
		binString = fmt.Sprintf("%s%.8b", binString, c)
	}
	return
}

func BinToString(s string) (binString string) {
	b := make([]byte, 0)
	for _, s := range strings.Fields(s) {
		n, _ := strconv.ParseUint(s, 2, 8)
		b = append(b, byte(n))
	}
	binString = string(b)
	return
}

func BinToStringOne(input string) (binString string) {
	r, _ := regexp.Compile("[0|1]{8}")

	match := r.FindAllString(input, -1)
	b := make([]byte, 0)
	for _, s := range match {
		n, _ := strconv.ParseUint(s, 2, 8)
		b = append(b, byte(n))
	}
	binString = string(b)
	return
}


func main() {
	a := makeTimestampT()

	fmt.Printf("%d \n", a)
	fmt.Println(dynamodbattribute.Marshal(1234))
	r := rand.New(rand.NewSource(55))
	ir := IntRange{1, 10}
	for i := 0; i < 10; i++ {
		fmt.Println(ir.NextRandom(r))
	}

	binaryString := stringToBin("Real Madrid")
	fmt.Println(binaryString)
	fmt.Println(BinToStringOne(binaryString))
}


func makeTimestampT() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
