package main

import (
	"fmt"
)

func main() {

	//Use the class factory to create an appliance of the requested type
	myAppliance, err := CreateAppliance(0)

	//if no errors start the appliance then print it's purpose
	if err == nil {
		myAppliance.Start()
		fmt.Println(myAppliance.GetPurpose())
	} else {
		//if error encountered, print the error
		fmt.Println(err)
	}

}
