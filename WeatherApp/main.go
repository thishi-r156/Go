package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/gorilla/mux"
    "io/ioutil"
    "os"
)

type APIKey struct{
    OpenWeatherMapAPI string `json:"OpenWeatherMapAPI"`
    LatLongAPI string `json:"LatLongAPI"`
}

type LocationResponse struct {
    Data []struct{
        Latitude float32 `json:"latitude"`
        Longitude float32 `json:"longitude`
    }
}

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int64 `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

var WeatherAPI, LatLongAPI string


func LoadAPIConfig(){
    configFile := ".apiConfig"
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config APIKey
	if err := json.Unmarshal(content, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

    WeatherAPI = config.OpenWeatherMapAPI
    LatLongAPI = config.LatLongAPI
}

func welcomeHome(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Hey User, Welcome to the Weather App\nGet instant Weather Updates of any place"))
}

func weatherReport(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    query := params["query"]
    requestURL := fmt.Sprintf("http://api.positionstack.com/v1/forward?access_key=%v&query=%v", LatLongAPI, query)
    res, err := http.Get(requestURL)
    if err != nil {
     log.Printf("error making http request: %s\n", err)
     return 
    }
    body, err := io.ReadAll(res.Body)
    if err != nil{
        log.Printf("Error reading response body: %s\n", err)
        return
    }

    var response LocationResponse
    err = json.Unmarshal(body, &response)
    if err != nil{
        log.Printf("Error in mapping response body to Json: %s\n", err)
        return
    }

    Lats := response.Data[0].Latitude
    Longs := response.Data[0].Longitude

    log.Println("Lattitude:", Lats)
    log.Println("Longitude:", Longs)
    weather_req := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=%v", Lats, Longs, WeatherAPI)
    weather_resp, err := http.Get(weather_req)
    if err != nil{
        log.Printf("error making http request: %s\n", err)
        return
    }
    weather_report, err1 := io.ReadAll(weather_resp.Body)
    if err1 != nil{
        log.Printf("Error reading response body: %s\n", err)
        return
    }

    var weather_response WeatherResponse
    err2 := json.Unmarshal(weather_report, &weather_response)
    if err2 != nil{
        log.Printf("Error in mapping response body to Json: %s\n", err)
        return
    }

    json.NewEncoder(w).Encode(&weather_response.Main)

}

func main(){
    LoadAPIConfig()
    r := mux.NewRouter()
    r.HandleFunc("/", welcomeHome)
    r.HandleFunc("/weather/{query}", weatherReport)
    log.Println("Starting on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))  
}