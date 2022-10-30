package main

import "time"

type Record struct {
	Timestamp     time.Time
	Temperature   float64
	Humidity      float64
	WindSpeed     float64
	Precipitation	float64
}

func ShouldWater(records []Record) (bool, error) {
	return true, nil
}
