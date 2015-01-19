package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type weatherData struct {
	Name    string `json:"name"`
	Weather []struct {
		ShortDescription string `json:"main"`
	} `json:"weather"`
	Main struct {
		Kelvin  float64 `json:"temp"`
		Celsius float64 `json:"temp_celsius"`
	} `json:"main"`
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	d.Main.Celsius = d.Main.Kelvin - 273.15

	return d, nil
}
