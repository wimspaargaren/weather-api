// Package main provides the weather API.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wimspaargaren/weather-api/pkg/weather"
)

func main() {
	weatherClient := weather.NewWeatherMapClient(os.Getenv("WEATHERMAP_API_KEY"))

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	weatherHandler := NewWeatherHandler(weatherClient)
	// Routes
	e.GET("/weather/current", weatherHandler.getCurrentWeather)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}

// WeatherHandler handles the weather API.
type WeatherHandler struct {
	weatherClient weather.Client
}

// NewWeatherHandler creates a new weather handler.
func NewWeatherHandler(weatherClient weather.Client) WeatherHandler {
	return WeatherHandler{
		weatherClient: weatherClient,
	}
}

// GetWeatherResponse is the response for the weather API get weather endpoint.
type GetWeatherResponse struct {
	Description string  `json:"description"`
	WindSpeed   float64 `json:"wind_speed"`
	Temperature float64 `json:"temperature"`
	Timestamp   int64   `json:"timestamp"`
}

func (w WeatherHandler) getCurrentWeather(c echo.Context) error {
	inputLocation := c.QueryParam("location")

	coordinate, err := w.weatherClient.GeoCoding(c.Request().Context(), inputLocation)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "sorry, we weren't able to find the location you were looking for")
	}

	currentWeather, err := w.weatherClient.CurrentWeather(c.Request().Context(), coordinate)
	if err != nil {
		log.Default().Println("error getting current weather:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "sorry, we weren't able to get the weather for the location you were looking for")
	}

	return c.JSON(http.StatusOK, GetWeatherResponse{
		Description: currentWeather.Description,
		WindSpeed:   currentWeather.WindSpeed,
		Temperature: currentWeather.Temperature,
		Timestamp:   currentWeather.Timestamp,
	})
}
