package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/njason/shouldwater/shouldwaterlib"
)

type Config struct {
	TomorrowioApiKey string `yaml:"tomorrowioApiKey"`
	MailChimp        struct {
		ApiKey     string `yaml:"apiKey"`
		TemplateId uint   `yaml:"templateId"`
		ListId     string `yaml:"listId"`
	}
	RecordsFile string `yaml:"recordsFile"`
	Lat         string `yaml:"lat"`
	Lng         string `yaml:"lng"`
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if len(os.Args) < 2 {
		historicalRecords, err := loadFreeRecords(config.RecordsFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if len(historicalRecords) < shouldwaterlib.HoursInWeek {
			log.Fatalln(errors.New("need at least a week's worth of data to run"))
		} else if len(historicalRecords) > shouldwaterlib.HoursInWeek {
			// trim to the last week of data
			historicalRecords = historicalRecords[len(historicalRecords)-shouldwaterlib.HoursInWeek:]
		}

		forecastRequest := NewTimelinesRequest(fmt.Sprintf("%s, %s", config.Lat, config.Lng), "metric", "1h", "nowPlus1h", "nowPlus5d")
		forecastRecordsRaw, err := DoTimelinesRequest(config.TomorrowioApiKey, forecastRequest)
		if err != nil {
			log.Fatalln(err.Error())
		}

		var forecastRecords []shouldwaterlib.Record
		for _, record := range forecastRecordsRaw.Data.Timelines[0].Intervals {
			forecastRecords = append(forecastRecords, shouldwaterlib.Record{
				Timestamp: record.StartTime,
				Temperature: record.Values.Temperature,
				Humidity: record.Values.Humidity,
				WindSpeed: record.Values.WindSpeed,
				Precipitation: record.Values.PrecipitationIntensity,
			})
		}

		shouldWater, err := shouldwaterlib.ShouldWater(historicalRecords, forecastRecords)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = archiveRecordsFile(config.RecordsFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if shouldWater {
			//err = createAndSendCampaign(config.MailChimp.ApiKey, config.MailChimp.TemplateId, config.MailChimp.ListId)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
	} else {
		switch os.Args[1] {
		case "record":
			err = RecordFreeTimelines(config.RecordsFile, config.Lat, config.Lng, config.TomorrowioApiKey)
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
