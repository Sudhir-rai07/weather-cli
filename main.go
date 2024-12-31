package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		Time      int64   `json:"last_updated_epoch"`
		Temp      float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		WindSpeed float64 `json:"wind_kph"`
		Humiduty  uint8   `json:"humidity"`
	} `json:"current"`
	Forecast struct {
		ForecastDay []struct {
			Time int64 `json:"date_epoch"`
			Day  struct {
				MaxT      float64 `json:"maxtemp_c"`
				MinT      float64 `json:"mintemp_c"`
				Wind      float64 `json:"maxwind_kph"`
				Humidity  uint8   `json:"avghumidity"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			}
		}
	} `json:"forecast"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load env")
	}
	ak := os.Getenv("API_KEY")
	city := ""
	if len(os.Args) < 2 {
		log.Fatal("Missing Argument - Area")
	}

	city = os.Args[1]

	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=3&aqi=no&alerts=no", ak, city)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Could not get weather data")
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("Bad Request")
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Could not get weather data")
		return
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Fatal("Faild to unmarshal data")
	}

	_, temp, text := weather.Location.Name, weather.Current.Temp, weather.Current.Condition.Text
	fmt.Println("TODAY")
	fmt.Printf("Temp : %.0fC\nCondition : %s\n", temp, text)

	fmt.Print("Forecast\n")
	fmt.Print("Day\tminT\t maxT\twind-kmph\ttext\thumidity\n")

	for _, fc := range weather.Forecast.ForecastDay {
		day := time.Unix(fc.Time, 0).Weekday().String()
		fmt.Printf("%s\t%.2fC\t%.2fC\t%.0f kmph\t%s\t%d\n",
			day,
			fc.Day.MinT,
			fc.Day.MaxT,
			fc.Day.Wind,
			fc.Day.Condition.Text,
			fc.Day.Humidity)
	}

}
