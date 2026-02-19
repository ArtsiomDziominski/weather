// Package main provides the core weather-fetching logic for the CLI app.
// It communicates with the OpenWeatherMap Current Weather Data API (v2.5).
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// baseURL is the OpenWeatherMap Current Weather endpoint.
const baseURL = "https://api.openweathermap.org/data/2.5/weather"

// --- JSON response structs ---

// weatherDesc holds a short weather description (e.g. "clear sky").
type weatherDesc struct {
	Description string `json:"description"`
}

// mainMetrics holds temperature, feels-like temperature, and humidity.
type mainMetrics struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity  int     `json:"humidity"`
}

// wind holds wind speed in m/s.
type wind struct {
	Speed float64 `json:"speed"`
}

// WeatherResponse is the top-level structure that maps the API JSON payload.
type WeatherResponse struct {
	Name    string        `json:"name"` // City name returned by the API
	Main    mainMetrics   `json:"main"`
	Weather []weatherDesc `json:"weather"` // Array; we use the first element
	Wind    wind          `json:"wind"`
	// The API returns a non-zero "cod" field (as a number or string) on errors.
	// We only need the numeric code for error detection.
	Cod     interface{} `json:"cod"`
	Message string      `json:"message"` // Error message from the API
}

// FetchWeather queries the OpenWeatherMap API for the current weather in city.
// It returns a populated *WeatherResponse on success, or an error on failure.
func FetchWeather(city, apiKey string) (*WeatherResponse, error) {
	// Build the request URL with query parameters.
	params := url.Values{}
	params.Set("q", city)
	params.Set("appid", apiKey)
	params.Set("units", "metric") // Celsius
	params.Set("lang", "en")

	requestURL := baseURL + "?" + params.Encode()

	// Make the HTTP GET request.
	resp, err := http.Get(requestURL) //nolint:gosec // URL is constructed from validated inputs
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON body regardless of status code — the API returns structured
	// error payloads (with "cod" and "message") even for 4xx responses.
	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// OpenWeatherMap signals errors via HTTP status codes (and a "message" field).
	if resp.StatusCode != http.StatusOK {
		if weather.Message != "" {
			return nil, fmt.Errorf("API error: %s", weather.Message)
		}
		return nil, fmt.Errorf("API returned unexpected status: %s", resp.Status)
	}

	// Sanity-check: we need at least one weather description element.
	if len(weather.Weather) == 0 {
		return nil, fmt.Errorf("API returned no weather data for %q", city)
	}

	return &weather, nil
}

// Display prints the weather data to stdout in a clean, human-readable format.
func Display(w *WeatherResponse) {
	// Capitalise the first letter of the description for nicer output.
	description := w.Weather[0].Description
	if len(description) > 0 {
		description = string(description[0]-32) + description[1:] // ASCII uppercase of first byte (safe for all-ASCII weather descriptions)
	}

	now := time.Now()
	fmt.Printf("City:        %s\n", w.Name)
	fmt.Printf("Time:        %s\n", now.Format("02 Jan 2006, 15:04"))
	fmt.Printf("Temperature: %.1f°C\n", w.Main.Temp)
	fmt.Printf("Feels like:  %.1f°C\n", w.Main.FeelsLike)
	fmt.Printf("Weather:     %s\n", description)
	fmt.Printf("Humidity:    %d%%\n", w.Main.Humidity)
	fmt.Printf("Wind speed:  %.1f m/s\n", w.Wind.Speed)
}
