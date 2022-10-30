package tomorrowio

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
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
		line := []string{
			interval.StartTime.String(),
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

func LoadFreeRecords(recordsFilename string) ([][]string, error) {
	file, err := os.Open(recordsFilename)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return data, nil
}
