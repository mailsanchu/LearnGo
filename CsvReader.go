package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Person struct {
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

func main() {
	csvFile, err := os.Open("people.csv")
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var people []Person
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		person := Person{
			FirstName: line[0],
			LastName:  line[1],
			Address: &Address{
				City:  line[2],
				State: line[3],
			},
		}
		fmt.Println(person.Address)
		people = append(people, person)
	}
	peopleJson, _ := json.MarshalIndent(people, "", "\t")
	fmt.Println(string(peopleJson))
}
