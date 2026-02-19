// Command weather is a CLI tool that fetches and displays current weather for
// a given city using the OpenWeatherMap API.
//
// Usage:
//
//	go run . [city]
//
// The WEATHER_API_KEY environment variable must be set to a valid
// OpenWeatherMap API key (or placed in a .env file in the same directory).
// If no city argument is provided, "London" is used as the default.
package main

import (
	"fmt"
	"os"
)

// defaultCity is used when the user does not supply a city argument.
const defaultCity = "London"

func main() {
	// --- Load .env file (if present) before reading env vars ---
	// This allows the user to store WEATHER_API_KEY in a .env file
	// instead of (or in addition to) setting it in the shell.
	if err := loadDotEnv(".env"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not read .env file: %v\n", err)
	}

	// --- Read API key from environment ---
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: WEATHER_API_KEY is not set.")
		fmt.Fprintln(os.Stderr, "Add it to the .env file in the project directory:")
		fmt.Fprintln(os.Stderr, "  WEATHER_API_KEY=your_api_key_here")
		fmt.Fprintln(os.Stderr, "Get a free key at https://openweathermap.org/api")
		os.Exit(1)
	}

	// --- Determine city from CLI arguments ---
	city := defaultCity
	if len(os.Args) > 1 {
		city = os.Args[1]
	}

	// --- Fetch weather data ---
	weather, err := FetchWeather(city, apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// --- Display the result ---
	Display(weather)
}
