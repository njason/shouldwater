package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/njason/shouldwater/shouldwater"
)

// RecordFreeTimelines reads in the maximum for free tier weather data into a csv file
func RecordFreeTimelines(filename string, lat string, lng string, tomorrowIoApiKey string) error {
	writeHeaders := false
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		writeHeaders = true
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if writeHeaders {
		err = writer.Write(append([]string{"startTime"}, timelineFields...))
		if err != nil {
			return err
		}
	}

	tomorrowIoRequest := NewTimelinesRequest(fmt.Sprintf("%s, %s", lat, lng), "metric", "1h", "nowMinus6h", "now")
	resp, err := DoTimelinesRequest(tomorrowIoRequest, tomorrowIoApiKey)
	if err != nil {
		return err
	}

	for _, interval := range resp.Data.Timelines[0].Intervals {
		timestamp, err := interval.StartTime.MarshalText()
		if err != nil {
			return err
		}

		line := []string{
			string(timestamp),
			fmt.Sprintf("%f", interval.Values.Temperature),
			fmt.Sprintf("%f", interval.Values.Humidity),
			fmt.Sprintf("%f", interval.Values.WindSpeed), 
			fmt.Sprintf("%f", interval.Values.PrecipitationIntensity)}

		err = writer.Write(line)
		if err != nil {
			return err
		}
	}

	writer.Flush()

	return nil
}

func loadFreeRecords(recordsFilename string) ([]shouldwater.Record, error) {
	file, err := os.Open(recordsFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var records []shouldwater.Record

	for i, row := range rows {
		if i == 0 {  // header row
			continue
		}

		var timestamp time.Time
		err := timestamp.UnmarshalText([]byte(row[0]))
		if err != nil {
			return nil, err
		}

		temperature, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, err
		}

		humidity, err :=  strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
		}

		windSpeed, err :=  strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, err
		}

		precipitation, err :=  strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, err
		}

		var record = shouldwater.Record{
			Timestamp: timestamp,
			Temperature: temperature,
			Humidity: humidity,
			WindSpeed: windSpeed,
			Precipitation: precipitation,
		}

		records = append(records, record)
	}

	return records, nil
}
