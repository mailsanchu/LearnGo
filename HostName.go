package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	name, _ := os.Hostname()
	springProfiles := os.Getenv("SPRING_PROFILES_ACTIVE")

	fmt.Println("hostname:", name)
	fmt.Println("Spring Profiles :", springProfiles)
	fmt.Println(runtime.NumCPU() + 1)
	fmt.Println(runtime.GOROOT())
}