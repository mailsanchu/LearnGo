package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func main() {
	viper.AutomaticEnv()

	for k, v := range viper.AllSettings() {
		fmt.Println(k, v)
	}
	fmt.Println(time.Now().Truncate(time.Second))

	dir, _ := os.Getwd()

	fmt.Println(dir)
	configBytes, err := ioutil.ReadFile(path.Join(dir, "devconfig.yaml"))
	if err != nil {
		fmt.Println(fmt.Errorf("IGNORING this Error when reading config: %v\n", err))
	}

	v1, err := readConfig("config", map[string]interface{}{
		"server.port":     8800,
	})

	if err != nil {
		panic(fmt.Errorf("Error when reading config: %v\n", err))
	}
	fmt.Println(len(configBytes))

	err = v1.MergeConfig(bytes.NewBuffer(configBytes)) // Find and read the config file
	if err != nil {                                    // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	port := v1.GetInt("server.port")
	hostname := v1.GetString("hostname")
	auth := v1.GetStringMapString("auth")

	fmt.Printf("Reading config for port = %d\n", port)
	fmt.Printf("Reading config for hostname = %s\n", hostname)
	fmt.Printf("Reading config for auth = %#v\n", auth)
}

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	for key, value := range defaults {
		v.Set(key, value)
	}
	return v, err
}

/*func init() {
	// Set viper path and read configuration
	if os.Getenv("ENV") == "PRODUCTION" {
		viper.AddConfigPath("config")
	} else {
		viper.AddConfigPath("devconfig")
	}
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		log.Fatalln("Fatal error config file", err)
	}
}*/
