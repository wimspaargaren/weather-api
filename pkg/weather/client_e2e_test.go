//go:build e2e

package weather

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeoCoding(t *testing.T) {
	weatherMapClient := NewWeatherMapClient(os.Getenv("WEATHERMAP_API_KEY"))
	coordinate, err := weatherMapClient.GeoCoding(context.Background(), "Schiphol")
	assert.NoError(t, err)
	assert.Equal(t, Coordinate{
		Lat: 52.3080392,
		Lon: 4.7621975,
	}, coordinate)
}

func TestCurrentWeather(t *testing.T) {
	weatherMapClient := NewWeatherMapClient(os.Getenv("WEATHERMAP_API_KEY"))
	currentWeather, err := weatherMapClient.CurrentWeather(context.Background(), Coordinate{
		Lat: 52.3080392,
		Lon: 4.7621975,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, currentWeather.Description)
	assert.NotZero(t, currentWeather.WindSpeed)
	assert.NotZero(t, currentWeather.Temperature)
	assert.NotZero(t, currentWeather.Timestamp)
}
