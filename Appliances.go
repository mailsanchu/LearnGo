package main

//import errors to log errors when they occur
import "errors"

//The main interface used to describe appliances
type Appliance interface{
	Start()
	GetPurpose() string
}

//Our appliance types
const (
	STOVE = iota
	FRIDGE
	//Now we support microwaves
	MICROWAVE
)

// define a stove struct, the struct contain a string representing the type name
type Stove struct{
	typeName string
}

//The stove struct implements the start() function
func (sv *Stove)Start(){
	sv.typeName = " Stove "
}

//The stove struct implements the GetPurpose() function
func (sv *Stove)GetPurpose() string{
	return "I am a " + sv.typeName + " I cook food!!"
}

// define a fridge struct, the struct contain a string representing the type name
type Fridge struct{
	typeName string
}

//The fridge struct implements the start() function
func (fr *Fridge)Start(){
	fr.typeName = " Fridge "
}

//The fridge struct implements the start() function
func (fr *Fridge)GetPurpose() string{
	return "I am a " + fr.typeName + " I cool stuff down!!"
}

type Microwave struct{
	typeName string
}

func (mr *Microwave)Start(){
	mr.typeName = " Microwave "
}

func (mr *Microwave)GetPurpose() string{
	return "I am a " + mr.typeName + " I heat stuff up!!"
}

//Function to create the appliances
func CreateAppliance(t int)(Appliance, error) {
	//Use a switch case to switch between types, if a type exist then error is nil (null)
	switch t{
	case STOVE:
		return new(Stove),nil
	case FRIDGE:
		return new(Fridge),nil
	case MICROWAVE:
		return new(Microwave),nil
	default:
		//if type is invalid, return an error
		return nil, errors.New("Invalid Appliance Type")
	}
}
