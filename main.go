package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

type weatherDataU struct {
	Name      string `json:"timezone"`
	Currently struct {
		Kelvin float64 `json:"temperature"`
	} `json:"currently"`
}

func main() {
	http.HandleFunc("/hi", hi)
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
	http.HandleFunc("/underground-weather/", func(w http.ResponseWriter, r *http.Request) {
		coordinates := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := queryU(coordinates)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hellooooo"))
}

func hi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi Andela"))
}

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=fa43ba207927450ab13764ecc3aa98b6&q=Nairobi" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

func queryU(coordinates string) (weatherDataU, error) {
	resp, err := http.Get("https://api.darksky.net/forecast/e07c745f1faa54174a07aae6324f92cb/" + coordinates)

	if err != nil {
		return weatherDataU{}, err
	}

	defer resp.Body.Close()

	var d weatherDataU

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherDataU{}, err
	}

	return d, nil
}
