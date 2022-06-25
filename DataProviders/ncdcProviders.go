package dataproviders

import (
	"encoding/json"
	"net/http"
	"fmt"
	"math"
	"os"
	"time"
	"log"
	"github.com/njason/shouldwater/models"
	"gopkg.in/yaml.v3"
)


const PrecipitationDataType = "PRCP"
const NCDCDateTimeFormat = "2006-01-02T15:04:05"

type NcdcConfig struct {
	Url string `yaml:"ncdcUrl"`
	Token string `yaml:"token"`
	DatasetId string `yaml:"datasetId"`
	Limit string `yaml:"limit"`
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

func loadConfig(configFile string) NcdcConfig {
	f, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
	}

	defer f.Close()

	var config NcdcConfig	
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
	}

	return config
}

func getQueryDateTimeFormat(time time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", time.Year(), time.Month(), time.Day())
}

func buildNcdcRequest(
	stationId string,
	daysToRequest int) *http.Request {

	config := loadConfig("config.yaml")

	today := time.Now()
	endDate := today.AddDate(0, 0, -1)
	startDate := endDate.AddDate(0, 0, -6)

	req, _ := http.NewRequest("GET", config.Url, nil)
	req.Header.Add("token", config.Token)
	query := req.URL.Query()
	query.Add("datasetid", config.DatasetId)
	query.Add("stationid", fmt.Sprintf("%s:%s", config.DatasetId, stationId))
	query.Add("startdate", getQueryDateTimeFormat(startDate))
	query.Add("enddate", getQueryDateTimeFormat(endDate))
	query.Add("limit", config.Limit)
	req.URL.RawQuery = query.Encode()

	return req
}

func getNCDCData(
	stationId string,
	daysToRequest int) NcdcResponse {

		req := buildNcdcRequest(stationId, daysToRequest)
	client := http.DefaultClient
	rawResp, _ := client.Do(req)
	defer rawResp.Body.Close()

	var resp NcdcResponse
	err := json.NewDecoder(rawResp.Body).Decode(&resp)
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
		log.Fatal()
	}

	return resp
}

// convertToInch will translate the NCDC value for precipitation, which is tenths of mm, to inches
func convertToInch(value int) float64 {
	return math.Round(float64(value)/10/25.4*100) / 100
}

func GetNCDCRainfall(
	stationId string,
	daysToRequest int) []models.RainfallDate {

	ncdcData := getNCDCData(stationId, daysToRequest)
	var rainfallPerDay [7]models.RainfallDate

	dayCounter := 0
	for _, result := range ncdcData.Results {
		if result.Datatype == PrecipitationDataType {
			var rainyDay models.RainfallDate
			parsedDateTime, err := time.Parse(NCDCDateTimeFormat, result.Date)
			if err != nil {
				fmt.Println(err)
			}
			rainyDay.Day = parsedDateTime
			rainyDay.RainfallInInches = convertToInch(result.Value)

			rainfallPerDay[dayCounter] = rainyDay
			dayCounter += 1
		}		
	}

	return rainfallPerDay[0:dayCounter]
}