package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"math"
	"time"
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
const WateringThreshold = 1  // in inches

// convertToInch will translate the NCDC value for precipitation, which is tenths of mm, to inches
func convertToInch(value int) float64 {
	return math.Round(float64(value) / 10 / 25.4 * 100) / 100
}

func getQueryFormat(time time.Time) string{
	return fmt.Sprintf("%d-%02d-%02d", time.Year(), time.Month(), time.Day())
}

func main() {
	token := os.Getenv("TOKEN")
	stationId := os.Getenv("STATIONID")

	client := http.DefaultClient

	end := time.Now()
	start := end.AddDate(0, 0, -7)

	req, _ := http.NewRequest("GET", "https://www.ncdc.noaa.gov/cdo-web/api/v2/data", nil)
	req.Header.Add("token", token)
	query := req.URL.Query()
	query.Add("datasetid", "GHCND")
	query.Add("stationid", fmt.Sprintf("GHCND:%s", stationId))
	query.Add("startdate", getQueryFormat(start))
	query.Add("enddate", getQueryFormat(end))
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

	for _, result := range resp.Results {
		if result.Datatype == PrecipitationDataType {
			totalPrecipRaw += result.Value
		}
	}

	totalPrecip := convertToInch(totalPrecipRaw)

	fmt.Printf("It's rained %.2f inches in the last week\n", totalPrecip)
	if (totalPrecip < WateringThreshold) {
		fmt.Println("You should water the trees.")
	} else {
		fmt.Println("No need to water the trees, it's rained enough.")
	}
}
