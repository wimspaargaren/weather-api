// Package weather provides a client for retrieving the weather.
package weather

import (
	"context"
	"fmt"
	"net/http"
)

// ErrNoCoordFoundForLocation is returned when no coordinates are found for the given location.
type ErrNoCoordFoundForLocation struct {
	Location string
}

func (e ErrNoCoordFoundForLocation) Error() string {
	return fmt.Sprintf("no coordinates found for location: %s", e.Location)
}

// ErrUnexpectedStatusCode is returned when the status code is not 200 OK.
type ErrUnexpectedStatusCode struct {
	StatusCode int
}

func (e ErrUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: %d", e.StatusCode)
}

// Weather represents the weather.
type Weather struct {
	Description string
	WindSpeed   float64
	Temperature float64
	Timestamp   int64
}

// Coordinate represents a coordinate.
type Coordinate struct {
	Lat float64
	Lon float64
}

// Client is the weather client.
type Client interface {
	GeoCoding(ctx context.Context, location string) (Coordinate, error)
	CurrentWeather(ctx context.Context, coordinate Coordinate) (*Weather, error)
}

// NewWeatherMapClient creates a client for the open weather map api.
func NewWeatherMapClient(apiKey string) Client {
	return &weatherMapClient{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
}
