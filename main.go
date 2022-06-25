package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/njason/shouldwater/dataproviders"
)

const WateringThreshold = 1  // in inches

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
	
	// Get NCDC Precipitation
	ncdcRainfall := dataproviders.GetNCDCRainfall(stationId, daysToRequest)	
	
	fmt.Println(ncdcRainfall)

	if len(ncdcRainfall) < 7 {
		fmt.Println("Insufficient data to total rain")
	} else {
		totalPrecip := 0.0
		for _, rainyDay := range ncdcRainfall {
			totalPrecip += rainyDay.RainfallInInches
		}
		
		fmt.Printf("It's rained %.2f inches in the last week\n", totalPrecip)
		if totalPrecip < WateringThreshold {
			fmt.Println("You should water the trees.")
		} else {
			fmt.Println("No need to water the trees, it's rained enough.")
		}
	}
}
