package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token string `yaml:"token"`
}

type NcdcResponse struct {
	Metadata struct {
		Resultset struct {
			Offset int `json:"offset"`
			Count  int `json:"count"`
			Limit  int `json:"limit"`
		} `json:"resultset"`
	} `json:"metadata"`
	Results []struct {
		Date       string `json:"date"`
		Datatype   string `json:"datatype"`
		Station    string `json:"station"`
		Attributes string `json:"attributes"`
		Value      int    `json:"value"`
	} `json:"results"`
}

const PrecipitationDataType = "PRCP"
const WateringThreshold = 1  // in inches

// convertToInch will translate the NCDC value for precipitation, which is tenths of mm, to inches
func convertToInch(value int) float64 {
	return math.Round(float64(value)/10/25.4*100) / 100
}

func getQueryFormat(time time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", time.Year(), time.Month(), time.Day())
}

func loadConfig() Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
	}

	return config
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Please provide station ID.")
		os.Exit(1)
	}

	config := loadConfig()
	stationId := flag.Arg(0)

	client := http.DefaultClient

	today := time.Now()
	end := today.AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -6)

	req, _ := http.NewRequest("GET", "https://www.ncdc.noaa.gov/cdo-web/api/v2/data", nil)
	req.Header.Add("token", config.Token)
	query := req.URL.Query()
	query.Add("datasetid", "GHCND")
	query.Add("stationid", fmt.Sprintf("GHCND:%s", stationId))
	query.Add("startdate", getQueryFormat(start))
	query.Add("enddate", getQueryFormat(end))
	query.Add("limit", "1000")
	req.URL.RawQuery = query.Encode()

	rawResp, _ := client.Do(req)
	defer rawResp.Body.Close()

	var resp NcdcResponse
	err := json.NewDecoder(rawResp.Body).Decode(&resp)
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
		return
	}

	totalPrecipRaw := 0
	dayCounter := 0

	for _, result := range resp.Results {
		if result.Datatype == PrecipitationDataType {
			totalPrecipRaw += result.Value
			dayCounter++
		}
	}

	if dayCounter < 7 {
		fmt.Println("Insufficient data to total rain")
	} else {
		totalPrecip := convertToInch(totalPrecipRaw)
		
		fmt.Printf("It's rained %.2f inches in the last week\n", totalPrecip)
		if totalPrecip < WateringThreshold {
			fmt.Println("You should water the trees.")
		} else {
			fmt.Println("No need to water the trees, it's rained enough.")
		}
	}
}
