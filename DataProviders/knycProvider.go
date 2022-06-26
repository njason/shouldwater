package dataproviders

import (
	"os"
	"io/ioutil"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"strings"
    "golang.org/x/net/html"
	"log"
	"github.com/njason/shouldwater/models"

	"gopkg.in/yaml.v3"
)

const TR_HEADER_TAG_COUNT = 7 // Number of tr tags to skip before data starts
const KNYCDateTimeFormat = "2006-01-02T15:04:05"

var _thTagCount = 0

type KNYCRow struct {
	Date string
	Time string
	Wind string
	Vis float64
	Weather string
	SkyCond string
	Temp struct {
		Air int
		Dewpoint int
		Max int
		Min int
	}
	RelativeHumidity int
	Windchill string
	HeatIndex string
	Pressure struct {
		Altimeter float64
		SeaLevel float64
	}
	Precipitation struct {
		LastHour float64
		Last3Hours float64
		Last6Hours float64
	}
}

type KNYCResponse struct {
	Data []KNYCRow
}

type KNYCConfig struct {
	Url string `yaml:"knycUrl"`
}

func loadKNYCConfig(configFile string) KNYCConfig {
	f, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
	}

	defer f.Close()

	var config KNYCConfig	
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err.Error(), http.StatusBadRequest)
	}

	return config
}

func isStartOfData(tkn html.Token) bool {
	if tkn.Data == "tr" {		
		if _thTagCount == TR_HEADER_TAG_COUNT {
			return true
		}
		_thTagCount++
	}

	return false
}

func readData(knycHtml *html.Tokenizer) string {
	for {
		tt := knycHtml.Next()
		switch tt {
			case html.EndTagToken:
				return ""
			case html.TextToken:
				data := knycHtml.Token().Data
				knycHtml.Next()
				knycHtml.Next()
				return data
		}
	}

	return ""
}

func readKNYCRow(knycHtml *html.Tokenizer) KNYCRow {
	var knycRow KNYCRow	
	knycHtml.Next()

	knycRow.Date = readData(knycHtml)			
	knycRow.Time = readData(knycHtml)
	knycRow.Wind = readData(knycHtml)
	knycRow.Vis, _ = strconv.ParseFloat(readData(knycHtml), 64)
	knycRow.Weather = readData(knycHtml)
	knycRow.SkyCond = readData(knycHtml)
	knycRow.Temp.Air, _ = strconv.Atoi(readData(knycHtml))
	knycRow.Temp.Dewpoint, _ = strconv.Atoi(readData(knycHtml))
	knycRow.Temp.Max, _ = strconv.Atoi(readData(knycHtml))
	knycRow.Temp.Min, _ = strconv.Atoi(readData(knycHtml))
	knycRow.RelativeHumidity, _ = strconv.Atoi(readData(knycHtml))
	knycRow.Windchill = readData(knycHtml)
	knycRow.HeatIndex = readData(knycHtml)
	knycRow.Pressure.Altimeter, _ = strconv.ParseFloat(readData(knycHtml), 64)
	knycRow.Pressure.SeaLevel, _ = strconv.ParseFloat(readData(knycHtml), 64)
	knycRow.Precipitation.LastHour, _ = strconv.ParseFloat(readData(knycHtml), 64)
	knycRow.Precipitation.Last3Hours, _ = strconv.ParseFloat(readData(knycHtml), 64)
	knycRow.Precipitation.Last6Hours, _ = strconv.ParseFloat(readData(knycHtml), 64)

	return knycRow
}

func readPrecipitationData(knycHtml *html.Tokenizer) KNYCResponse {

	var knycResponse KNYCResponse
	var rows [100]KNYCRow
	rowCounter := 0

	for {
		tt := knycHtml.Next()
	
		switch tt {
			case html.ErrorToken :
				knycResponse.Data = rows[0:rowCounter-2]
				return knycResponse
			case html.StartTagToken:
				if (knycHtml.Token().Data == "tr") {
					rows[rowCounter] = readKNYCRow(knycHtml)
					rowCounter++
				}
		}
	}

	return knycResponse
}

func parseHtml(htmlStr string) KNYCResponse {
    knycHtml := html.NewTokenizer(strings.NewReader(htmlStr))

	for {
		tt := knycHtml.Next()
		switch tt{
			case html.ErrorToken :
				var k KNYCResponse
				return k
	
			case html.StartTagToken :
				if isStartOfData(knycHtml.Token()) {
					fmt.Println("Start Reading Data")
					precipitationData := readPrecipitationData(knycHtml)
					fmt.Println("Done Reading Data")
					return precipitationData
				}				
		}
	}

}

func getKNYCData() KNYCResponse {
	knycConfig := loadKNYCConfig("config.yaml")
	req, _ := http.NewRequest("GET", knycConfig.Url, nil)
	client := http.DefaultClient
	rawResp, _ := client.Do(req)

	htmlStr, err := ioutil.ReadAll(rawResp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	knycResponse := parseHtml(string(htmlStr))
	return knycResponse
}

func GetKNYCRainfall() []models.RainfallDate {

	knycData := getKNYCData()

	var knycRainfall [7]models.RainfallDate
	dayCounter := 0
	dailyRainfallMap := make(map[string]float64)
	for _,knycRow := range knycData.Data {

		year, month, _ := time.Now().Date()
		rowDate := fmt.Sprintf("%4d-%02d-%02sT00:00:00", year, month, knycRow.Date)

		dailyRainfallMap[rowDate] += knycRow.Precipitation.LastHour 
		dailyRainfallMap[rowDate] += knycRow.Precipitation.Last3Hours 
		dailyRainfallMap[rowDate] += knycRow.Precipitation.Last6Hours
	} 

	for dateStr, rainfall := range dailyRainfallMap {
		var rainfallDate models.RainfallDate
		rainfallDate.Day, _ = time.Parse(KNYCDateTimeFormat, dateStr)
		rainfallDate.RainfallInInches = rainfall
		knycRainfall[dayCounter] = rainfallDate

		dayCounter++
	}

	return knycRainfall[0:dayCounter]
}