package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// WeatherResponse defines the structure of the JSON response

type WeatherResponse struct {
	City         string  `json:"city"`
	Date         string  `json:"date"`
	TemperatureF float64 `json:"temperature_f"`
}

func getMockWeatherData(city, date string) float64 {

	return 75.5
}

// logRequestMiddleware that logs each and every incoming request.

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
		//fmt.Printf("Headers: %v\n", r.Header)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// weatherHandler handles the GET request for the weather data
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	date := r.URL.Query().Get("date")

	// validation
	if city == "" {
		http.Error(w, "Missing required parameter: city", http.StatusBadRequest)
		return
	}
	if date == "" {
		http.Error(w, "Missing required parameter: date", http.StatusBadRequest)
	}

	temperature := getMockWeatherData(city, date)

	//create response
	response := WeatherResponse{
		City:         city,
		Date:         date,
		TemperatureF: temperature,
	}
	// set response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode the response to JSON and send it

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func main() {

	// Create a new Chi router
	r := chi.NewRouter()

	//Adding middleware for logging and recovering from panics.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logRequestMiddleware)

	// set up /v1/weather endpoint
	r.Get("/v1/weather/{city}", weatherHandler)

	// start the http server
	fmt.Println("Starting server on : 8080....")
	log.Fatal(http.ListenAndServe(":8080", r))
}
