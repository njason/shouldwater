package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/njason/shouldwater/tomorrowio"
)

type Config struct {
	TomorrowIoApiKey string `yaml:"tomorrowioApiKey"`
	RecordsFile	string `yaml:"recordsFile"`
	Lat string `yaml:"lat"`
	Lng string `yaml:"lng"`
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if len(os.Args) < 2 {
		records, err := loadRecords(config.RecordsFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		shouldWater, err := ShouldWater(records)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if shouldWater {
			log.Println("Should water")
		}
	} else {
		switch os.Args[1] {
		case "record":
			err = tomorrowio.RecordFreeTimelines(config.RecordsFile, config.Lat, config.Lng, config.TomorrowIoApiKey)
			if err != nil {
				log.Fatalln(err.Error())
			}
		default:
			log.Fatalf("unknown command '%s'\n", os.Args[1])
		}
	}
}

func loadConfig() (Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func loadRecords(recordsFilename string) ([]Record, error) {
	rawRecords, err := tomorrowio.LoadFreeRecords(recordsFilename)
	if err != nil {
		return []Record{}, err
	}

	var records []Record
	for i, line := range rawRecords {
		if i > 0 {  // omit header line
			timestamp, err := time.Parse("D", line[0])
			if err != nil {
				return []Record{}, err
			}

			temperature, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				return []Record{}, err
			}

			humidity, err :=  strconv.ParseFloat(line[2], 64)
			if err != nil {
				return []Record{}, err
			}

			windSpeed, err :=  strconv.ParseFloat(line[2], 64)
			if err != nil {
				return []Record{}, err
			}

			precipitation, err :=  strconv.ParseFloat(line[3], 64)
			if err != nil {
				return []Record{}, err
			}

			var record = Record{
				Timestamp: timestamp,
				Temperature: temperature,
				Humidity: humidity,
				WindSpeed: windSpeed,
				Precipitation: precipitation,
			}

			records = append(records, record)
		}
	}

	return records, nil
}
