package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	token := os.Getenv("TOKEN")

	client := http.DefaultClient

	req, _ := http.NewRequest("GET", "https://www.ncdc.noaa.gov/cdo-web/api/v2/data?datasetid=GHCND&stationid=GHCND:USW00094728&startdate=2022-01-16&enddate=2022-01-22", nil)
	req.Header.Add("token", token)

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}
