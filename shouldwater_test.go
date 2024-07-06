package shouldwater

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// generateWeatherRecords generates WeatherRecord instances for the amount of hours
// For now, all the timestamps are the same. This should change if hourly timestamp sequencess are validated
func generateWeatherRecords(hourAmount int, temperature float64, precipitation float64) []WeatherRecord {
	records := make([]WeatherRecord, 0, hourAmount)
	start := time.Now()//.Add(backForward * time.Duration(hourAmount) * time.Hour)  // Go back 168 hours

	for i := 0; i < hourAmount; i++ {
		timestamp := start//start.Add(time.Duration(i) * time.Hour)

		// Predefined weather conditions with zero precipitation
		weatherRecord := WeatherRecord{
			Timestamp: timestamp,
			Temperature:   temperature,
			Humidity:      0.0,
			WindSpeed:     0.0,
			Precipitation: precipitation,
		}
		records = append(records, weatherRecord)
	}

	return records
}


func TestNotEnoughData(t *testing.T) {
	_, err := ShouldWater([]WeatherRecord{}, []WeatherRecord{})

	require.NotNil(t, err)
}

func TestMaxWatering(t *testing.T) {
	historicalRecords := generateWeatherRecords(HoursInWeek, 25, 0)
	forecastRecords := generateWeatherRecords(HoursInFiveDays, 25, 0)

	amountToWater, err := ShouldWater(historicalRecords, forecastRecords)

	require.Nil(t, err)
	require.Equal(t, WateringMax, amountToWater)
}

func TestMinWatering(t *testing.T) {
	historicalRecords := generateWeatherRecords(HoursInWeek, 25, 25.3)
	forecastRecords := generateWeatherRecords(HoursInFiveDays, 25, 25.3)

	amountToWater, err := ShouldWater(historicalRecords, forecastRecords)

	require.Nil(t, err)
	require.Equal(t, 0.0, amountToWater)
}

func TestOneThirdWatering(t *testing.T) {
	// 2/3 of max precipitation divided by sum of historical and forecase hours
	historicalRecords := generateWeatherRecords(HoursInWeek, 25, 0.07643518518)
	forecastRecords := generateWeatherRecords(HoursInFiveDays, 25, 0.07643518518)

	amountToWater, err := ShouldWater(historicalRecords, forecastRecords)

	require.Nil(t, err)
	require.Equal(t, math.Round(WateringMax/3*100000)/100000, math.Round(amountToWater*100000)/100000)
}
