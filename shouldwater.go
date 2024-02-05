package shouldwater

import (
	"fmt"
	"time"
)

type WeatherRecord struct {
	Timestamp     time.Time
	Temperature   float64
	Humidity      float64
	WindSpeed     float64
	Precipitation float64
}

const HoursInWeek = 7 * 24
const HoursInFiveDays = 5 * 24

const HighTempHistoricalPrecipitationMax = 25.4 // 1 inch in mm
const HighTempForecastPrecipitationMax = 25.4   // 1 inch in mm
const HistoricalPrecipitationMax = 20.32        // .8 inches in mm
const ForecastPrecipitationMax = 12.7           // .5 inches in mm
const WateringMax = 75.71                       // liters

// ShouldWater returns the amount of liters required for watering an unestablished street tree
// given a weeks worth of historical and five days worth of forecasted weather data, in hourly granularity
func ShouldWater(
	historicalRecords []WeatherRecord,
	forecastRecords []WeatherRecord,
) (float64, error) {

	err := validateShouldWaterInputData(historicalRecords, forecastRecords)
	if err != nil {
		return 0.0, err
	}

	totalHistoricalPrecipitation := totalNonFastFallPrecipitation(historicalRecords)
	totalForecastPrecipitation := totalNonFastFallPrecipitation(forecastRecords)
	totalPrecipitation := totalHistoricalPrecipitation + totalForecastPrecipitation
	averageHistoricalHighTemperature := averageDayHighTemperature(historicalRecords)

	var totalPrecipitationMax float64
	if averageHistoricalHighTemperature > 29.4 { // 85 F in C
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

func ShouldHaveWatered(weatherRecords []WeatherRecord) (float64, error) {

	err := validateShouldHaveWateredInputData(weatherRecords)
	if err != nil {
		return 0.0, err
	}

	totalPrecipitation := totalNonFastFallPrecipitation(weatherRecords)
	averageHistoricalHighTemperature := averageDayHighTemperature(weatherRecords)

	var totalPrecipitationMax float64
	if averageHistoricalHighTemperature > 29.4 { // 85 F in C
		totalPrecipitationMax = HighTempHistoricalPrecipitationMax * 2
	} else {
		totalPrecipitationMax = HistoricalPrecipitationMax * 2
	}

	if totalPrecipitation < totalPrecipitationMax {
		percentPrecipitation := totalPrecipitation / totalPrecipitationMax
		return WateringMax * percentPrecipitation, nil
	}

	return 0.0, nil
}

func validateShouldWaterInputData(
	historicalRecords []WeatherRecord,
	forecastRecords []WeatherRecord,
) error {
	if len(historicalRecords) != HoursInWeek {
		return fmt.Errorf("need exactly a week's worth of historical data to run (%d records)", HoursInWeek)
	}

	if len(forecastRecords) != HoursInFiveDays {
		return fmt.Errorf("need exactly five days worth of forecast data to run (%d records)", HoursInFiveDays)
	}

	return nil
}

func validateShouldHaveWateredInputData(weatherRecords []WeatherRecord) error {
	if len(weatherRecords) != HoursInWeek * 2 {
		return fmt.Errorf("need exactly two week's worth of hourly records (%d records)", HoursInWeek * 2)
	}

	return nil
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
