package main

import (
	"fmt"
	"os"
	"log"
	//"flag"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TomorrowIoApiKey string `yaml:"tomorrowioApiKey"`
}

const WateringThreshold = 1  // in inches

func loadConfig() Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return config
}

func main() {
	/*flag.Parse()
	if flag.NArg() > 0 {
		fmt.Println("Please provide station ID.")
		os.Exit(1)
	}

	reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
	placeUrl := reportCmd.String("place", "", "Place URL, ex. jersey-city")
	category := reportCmd.String("category", "", "ex. trees")

	if len(os.Args) < 2 {
		fmt.Println("expected 'report' subcommand")
		os.Exit(1)
	}*/

	config := loadConfig()

	result, err := tomorrowIoRequest(config.TomorrowIoApiKey)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(result)
}
