package shouldwaterlib

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

const HoursInWeek = 7 * 24
const HoursInFiveDays = 5 * 24

func ShouldWater(historicalRecords []Record, forecastRecords []Record) (bool, error) {
	if len(historicalRecords) != HoursInWeek {
		return false, errors.New("need exactly a week's worth of historical data to run")
	}

	if len(forecastRecords) != HoursInFiveDays {
		return false, errors.New("need exactly five days worth of forecast data to run")
	}

	totalHistoricalPrecipitation := totalPrecipitation(historicalRecords)

	if totalHistoricalPrecipitation < 25.4 {  // 1 inch in mm
		
		totalForecastPrecipitation := totalPrecipitation(forecastRecords)
		if totalForecastPrecipitation < 50.8 {  // 2 inches in mm
			return true, nil
		}
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
