// Package app bootstraps the weather app.
package app

import (
	"net/http"
	"os"

	internalHTTP "github.com/wimspaargaren/weather-api/internal/http"
	"github.com/wimspaargaren/weather-api/pkg/api"
	"github.com/wimspaargaren/weather-api/pkg/weather"
)

// Run starts the weather app.
func Run() error {
	weatherClient := weather.NewWeatherMapClient(os.Getenv("WEATHERMAP_API_KEY"))
	weatherHandler := internalHTTP.NewWeatherHandler(weatherClient)

	server := api.NewServer()
	server.RegisterRoute(http.MethodGet, "/weather/current", weatherHandler.GetCurrentWeather)
	return server.Run()
}
