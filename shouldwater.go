package shouldwater

import (
	"errors"
	"time"

	"golang.org/x/tools/go/analysis/passes/nilness"
)

type WeatherRecord struct {
	Timestamp     time.Time
	Temperature   float64
	Humidity      float64
	WindSpeed     float64
	Precipitation	float64
}

const HoursInWeek = 7 * 24
const HoursInFiveDays = 5 * 24

const HighTempHistoricalPrecipitationMax = 25.4  // 1 inch in mm
const HighTempForecastPrecipitationMax = 25.4  // 1 inch in mm
const HistoricalPrecipitationMax = 20.32  // .8 inches in mm
const ForecastPrecipitationMax = 12.7  // .5 inches in mm
const WateringMax = 75.71  // liters

func ShouldWater(historicalRecords []WeatherRecord, forecastRecords []WeatherRecord) (float64, error) {
	if len(historicalRecords) != HoursInWeek {
		return 0.0, errors.New("need exactly a week's worth of historical data to run")
	}

	if len(forecastRecords) != HoursInFiveDays {
		return 0.0, errors.New("need exactly five days worth of forecast data to run")
	}

	totalHistoricalPrecipitation := totalNonFastFallPrecipitation(historicalRecords)
	averageHistoricalHighTemperature := averageDayHighTemperature(historicalRecords)
	totalForecastPrecipitation := totalNonFastFallPrecipitation(forecastRecords)

	totalPrecipitation := totalHistoricalPrecipitation + totalForecastPrecipitation

	var totalPrecipitationMax float64
	if averageHistoricalHighTemperature > 29.4 {  // 85 F in C
		totalPrecipitationMax = HighTempHistoricalPrecipitationMax + HighTempForecastPrecipitationMax
	} else {
		totalPrecipitationMax = HistoricalPrecipitationMax + ForecastPrecipitationMax
	}

	if totalPrecipitation < totalPrecipitationMax {
		percentPrecipitation := totalPrecipitation / totalPrecipitationMax
		return WateringMax * percentPrecipitation, nil
	}

	return 0.0, nil
}

func totalNonFastFallPrecipitation(records []WeatherRecord) float64 {
	var total float64 
	for _, record := range records {
		if record.Precipitation < 25.4 {  // 1 inches in mm
			total += record.Precipitation
		}
	}

	return total
}

func averageDayHighTemperature(records []WeatherRecord) float64 {
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
