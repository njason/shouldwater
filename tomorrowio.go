package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type tomorrowIoResponse struct {
	Data struct {
		Timelines []struct {
			Timestep  string    `json:"timestep"`
			EndTime   time.Time `json:"endTime"`
			StartTime time.Time `json:"startTime"`
			Intervals []struct {
				StartTime time.Time `json:"startTime"`
				Values    struct {
					PrecipitationIntensity float64 `json:"precipitationIntensity"`
				} `json:"values"`
			} `json:"intervals"`
		} `json:"timelines"`
	} `json:"data"`
	Code    *int   `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func tomorrowIoRequest(apiKey string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.tomorrow.io/v4/timelines", nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()
	query.Add("location", "40.7831, -73.9713")
	query.Add("fields", "precipitationIntensity")
	query.Add("units", "imperial")
	query.Add("timesteps", "1h")
	query.Add("startTime", "nowMinus6h")
	query.Add("endTime", "now")
	query.Add("apikey", apiKey)
	req.URL.RawQuery = query.Encode()

	req.Header.Add("accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	zippedBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	gzreader, err := unzip(zippedBody)
	if err != nil {
		return "", err
	}

	var resp tomorrowIoResponse
	err = json.NewDecoder(gzreader).Decode(&resp)
	if err != nil {
		return "", err
	}

	if resp.Code != nil {
		return "", errors.New(resp.Message)
	}

	return fmt.Sprintf("%.2f inches", totalPrecipitation(resp)), nil
}

func unzip(zipped []byte) (io.Reader, error) {
	reader := bytes.NewReader([]byte(zipped))

	return gzip.NewReader(reader)
}

func totalPrecipitation(resp tomorrowIoResponse) float64 {
	var total = .0

	for _, interval := range resp.Data.Timelines[0].Intervals {
		total += interval.Values.PrecipitationIntensity
	}

	return total
}
