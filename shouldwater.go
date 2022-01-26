package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"math"
)


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

// convertToInch will translate the NCDC value for precipitation, which is tenths of mm, to inches
func convertToInch(value int) float64 {
	return math.Round(float64(value) / 10 / 25.4 * 100) / 100
}

func main() {
	token := os.Getenv("TOKEN")

	client := http.DefaultClient

	req, _ := http.NewRequest("GET", "https://www.ncdc.noaa.gov/cdo-web/api/v2/data?datasetid=GHCND&stationid=GHCND:USW00094728&startdate=2022-01-16&enddate=2022-01-22", nil)
	req.Header.Add("token", token)

	rawResp, _ := client.Do(req)
	defer rawResp.Body.Close()

	var resp NcdcResponse
	err := json.NewDecoder(rawResp.Body).Decode(&resp)
	if err != nil {
        fmt.Println(err.Error(), http.StatusBadRequest)
        return
    }

	totalPrecip := 0

	for _, result := range resp.Results {
		if result.Datatype == PrecipitationDataType {
			totalPrecip += result.Value
		}
	}

	fmt.Println(convertToInch(totalPrecip))
}
