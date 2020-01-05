package config

import "fmt"

// Configurations exported
type Configurations struct {
	Server       ServerConfigurations
	Database     DatabaseConfigurations
	DME2_GRM_SERVER_PORT  string
	EXAMPLE_VAR  string
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port int
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBName     string
	DBUser     string
	DBPassword string
}

func main()  {
	fmt.Println("Sanchu")
}
