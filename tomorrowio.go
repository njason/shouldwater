package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: 60 * time.Second,
}

var timelineFields = []string{"temperature", "humidity", "windSpeed", "precipitationIntensity"}

type timelinesRequest struct {
	location	string
	fields	[]string
	units	string
	timesteps string
	startTime	string
	endTime	string
}

type timelinesResponse struct {
	Data struct {
		Timelines []struct {
			Timestep  string    `json:"timestep"`
			EndTime   time.Time `json:"endTime"`
			StartTime time.Time `json:"startTime"`
			Intervals []struct {
				StartTime time.Time `json:"startTime"`
				Values    struct {  // https://docs.tomorrow.io/reference/data-layers-core
					Temperature float64 `json:"temperature"`
					Humidity float64 `json:"humidity"`
					WindSpeed float64 `json:"windSpeed"`
					PrecipitationIntensity float64 `json:"precipitationIntensity"`
				} `json:"values"`
			} `json:"intervals"`
		} `json:"timelines"`
	} `json:"data"`
	Code    *int   `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewTimelinesRequest(location string, units string, timesteps string,
	startTime string, endTime string) timelinesRequest {
		return timelinesRequest{
			location: location,
			fields: timelineFields,
			units: units,
			timesteps: timesteps,
			startTime: startTime,
			endTime: endTime,
		}
}

func DoTimelinesRequest(request timelinesRequest, apiKey string) (timelinesResponse, error) {
	req, err := http.NewRequest("GET", "https://api.tomorrow.io/v4/timelines", nil)
	if err != nil {
		return timelinesResponse{}, err
	}

	query := req.URL.Query()
	query.Add("location", request.location)
	query.Add("units", request.units)
	query.Add("timesteps", request.timesteps)
	query.Add("startTime", request.startTime)
	query.Add("endTime", request.endTime)
	query.Add("apikey", apiKey)

	for _, field := range request.fields {
		query.Add("fields", field)
	}
	req.URL.RawQuery = query.Encode()

	req.Header.Add("accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")

	res, err := client.Do(req)
	if err != nil {
		return timelinesResponse{}, err
	}
	defer res.Body.Close()

	zippedBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return timelinesResponse{}, err
	}

	gzreader, err := unzip(zippedBody)
	if err != nil {
		return timelinesResponse{}, err
	}

	var resp timelinesResponse
	err = json.NewDecoder(gzreader).Decode(&resp)
	if err != nil {
		return timelinesResponse{}, err
	}

	if resp.Code != nil {
		return timelinesResponse{}, errors.New(resp.Message)
	}

	return resp, nil
}

func unzip(zipped []byte) (io.Reader, error) {
	reader := bytes.NewReader([]byte(zipped))

	return gzip.NewReader(reader)
}
