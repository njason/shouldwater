package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"github.com/njason/shouldwater/dataproviders"
	"github.com/njason/shouldwater/models"
)

const WateringThreshold = 1  // in inches

func aggregateRainfallData(
	rainfallDays []models.RainfallDate, 
	totalRainfallMap map[time.Time]float64) {
	for _, rainfallDay := range rainfallDays {
		totalRainfallMap[rainfallDay.Day] += rainfallDay.RainfallInInches
	}	
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Please provide station ID.")
		os.Exit(1)
	}

	fmt.Println("Done loading config.");
	stationId := flag.Arg(0)

	fmt.Println("StationId: ", stationId);
	daysToRequest := 7
	
	totalRainfallMap := make(map[time.Time]float64)
	// Get NCDC Precipitation
	fmt.Println("Start NCDC Data")
	ncdcRainfall := dataproviders.GetNCDCRainfall(stationId, daysToRequest)	
	fmt.Println(ncdcRainfall)
	aggregateRainfallData(ncdcRainfall, totalRainfallMap)

	fmt.Println("Start KNYC Data")
	knycRainfall := dataproviders.GetKNYCRainfall()
	fmt.Println(knycRainfall)
	aggregateRainfallData(knycRainfall, totalRainfallMap)

	if len(totalRainfallMap) < 7 {
		fmt.Println("Insufficient data to total rain. Only received", len(ncdcRainfall), "days of data.")
	} else {
		totalPrecip := 0.0
		for _, rainfall := range totalRainfallMap {
			totalPrecip += rainfall
		}
		
		fmt.Printf("It's rained %.2f inches in the last week\n", totalPrecip)
		if totalPrecip < WateringThreshold {
			fmt.Println("You should water the trees.")
		} else {
			fmt.Println("No need to water the trees, it's rained enough.")
		}
	}
}
