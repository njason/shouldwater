package shouldwater

import (
	"errors"
	"time"
)

type Record struct {
	Timestamp     time.Time
	Temperature   float64
	Humidity      float64
	WindSpeed     float64
	Precipitation	float64
}

func ShouldWater(records []Record) (bool, error) {
	if len(records) != 168 {  // 7 days worth of hourly data
		return false, errors.New("need exactly a week's worth of data to decide")
	}

	totalPrecipitation := totalPrecipitation(records)

	if totalPrecipitation > 25.4 {  // 1 inch in mm
		return true, nil
	}

	return false, nil
}

func totalPrecipitation(records []Record) float64 {
	var total float64 
	for _, record := range records {
		total += record.Precipitation
	}

	return total
}
