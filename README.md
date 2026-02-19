# Weather CLI

A simple, production-quality CLI application written in Go that fetches and displays the current weather for any city using the [OpenWeatherMap](https://openweathermap.org/) free API.

## Features

- Displays temperature (°C), feels-like, weather description, humidity, and wind speed
- Defaults to **London** if no city is specified
- Friendly error messages for missing API key, invalid city, or network issues
- Zero external dependencies — standard library only

## Getting a Free API Key

1. Sign up for free at <https://openweathermap.org/api>
2. Navigate to **API keys** in your account dashboard
3. Copy your default key (it becomes active within a few minutes)

## Setting the Environment Variable

**Windows (PowerShell)**

```powershell
$env:WEATHER_API_KEY = "your_api_key_here"
```

**Windows (Command Prompt)**

```cmd
set WEATHER_API_KEY=your_api_key_here
```

**Linux / macOS**

```bash
export WEATHER_API_KEY=your_api_key_here
```

## Running the App

Make sure you are in the project directory:

```powershell
cd c:\Users\hitma\GolandProjects\weather
```

**Run directly with `go run`:**

```powershell
# Default city (London)
go run .

# Specific city
go run . Tbilisi
go run . "New York"
```

**Or build a binary first:**

```powershell
go build -o weather.exe .
.\weather.exe Tbilisi
```

## Example Output

```
City:        Tbilisi
Temperature: 18.0°C
Feels like:  16.2°C
Weather:     Clear sky
Humidity:    60%
Wind speed:  3.5 m/s
```

## Error Examples

```
# Missing API key
Error: WEATHER_API_KEY environment variable is not set.

# Invalid city
Error: API error: city not found

# Network issue
Error: network error: <details>
```

## Project Structure

```
weather/
├── main.go      # CLI entry point — arg parsing, env var, error handling
├── weather.go   # API client, JSON structs, display logic
├── go.mod       # Go module definition (standard library only)
└── README.md
```