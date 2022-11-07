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

	totalHistoricalPrecipitation := totalNonFastFallPrecipitation(historicalRecords)
	averageHistoricalHighTemperature := averageDayHighTemperature(historicalRecords)
	totalForecastPrecipitation := totalNonFastFallPrecipitation(forecastRecords)

	if averageHistoricalHighTemperature > 29.4 {  // 85 F in C
		if totalHistoricalPrecipitation < 25.4 {  // 1 inch in mm
			if totalForecastPrecipitation < 25.4 {  // 1 inches in mm
				return true, nil
			}
		}
	} else {
		if totalHistoricalPrecipitation < 20.32 {  // .8 inches in mm
			if totalForecastPrecipitation < 12.7 {  // .5 inches in mm
				return true, nil
			}
		}
	}


	return false, nil
}

func totalNonFastFallPrecipitation(records []Record) float64 {
	var total float64 
	for _, record := range records {
		if record.Precipitation < 25.4 {  // 1 inches in mm
			total += record.Precipitation
		}
	}

	return total
}

func averageDayHighTemperature(records []Record) float64 {
	dayHighs := make(map[string]float64)

	for _, record := range records {
		day := record.Timestamp.Format("YYYYMMDD")

		if dayHigh, ok := dayHighs[day]; ok {
			if dayHigh < record.Temperature {
				dayHighs[day] = record.Temperature
			}
		} else {
			dayHighs[day] = record.Temperature
		}
	}

	var dayHighSum float64
	for _, dayHigh := range dayHighs {
        dayHighSum += dayHigh
    }

	return dayHighSum / float64(len(dayHighs))
}
