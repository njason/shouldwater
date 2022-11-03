package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/njason/shouldwater/shouldwater"
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
		records, err := loadFreeRecords(config.RecordsFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		shouldWater, err := shouldwater.ShouldWater(records)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = archiveRecordsFile(config.RecordsFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if shouldWater {
			log.Println("Should water")
		}
	} else {
		switch os.Args[1] {
		case "record":
			err = RecordFreeTimelines(config.RecordsFile, config.Lat, config.Lng, config.TomorrowIoApiKey)
			if err != nil {
				log.Fatalln(err.Error())
			}
		default:
			log.Fatalf("unknown command '%s'. The only command is 'record'\n", os.Args[1])
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

func archiveRecordsFile(recordsFile string) error {
	fileName := strings.TrimSuffix(recordsFile, filepath.Ext(recordsFile))
	archiveFile := fmt.Sprintf("%s_archive.csv", fileName)
	err := os.Rename(recordsFile, archiveFile)

	if err != nil {
		return err
	}

	return nil
}
