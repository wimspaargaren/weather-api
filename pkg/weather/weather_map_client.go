package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
)

// weatherMapClient is a client for the open weather map api.
type weatherMapClient struct {
	apiKey     string
	httpClient *http.Client
}

const (
	weatherMapBaseURL = "api.openweathermap.org"
	geoV1Path         = "/geo/1.0/direct"
	weatherV25Path    = "/data/2.5/weather"
)

type geoCodingWeatherMapResponse []struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

// GeoCoding returns the coordinate for the given location.
func (c *weatherMapClient) GeoCoding(ctx context.Context, location string) (Coordinate, error) {
	req, err := c.newWeatherMapGetRequest(ctx, geoV1Path)
	if err != nil {
		return Coordinate{}, err
	}
	queryParams := req.URL.Query()
	queryParams.Set("q", location)
	queryParams.Set("limit", "1")
	req.URL.RawQuery = queryParams.Encode()

	weatherMapResponse := geoCodingWeatherMapResponse{}
	err = c.executeRequest(req, &weatherMapResponse)
	if err != nil {
		return Coordinate{}, err
	}

	if len(weatherMapResponse) == 0 {
		return Coordinate{}, ErrNoCoordFoundForLocation{
			Location: location,
		}
	}
	return Coordinate{
		Lat: weatherMapResponse[0].Lat,
		Lon: weatherMapResponse[0].Lon,
	}, nil
}

type currentWeatherMapResponse struct {
	Coord      weatherMapCoord     `json:"coord"`
	Weather    []weatherMapWeather `json:"weather"`
	Base       string              `json:"base"`
	Main       weatherMapMain      `json:"main"`
	Visibility int                 `json:"visibility"`
	Wind       weatherMapWind      `json:"wind"`
	Clouds     weatherMapClouds    `json:"clouds"`
	Dt         int                 `json:"dt"`
	Sys        weatherMapSys       `json:"sys"`
	Timezone   int                 `json:"timezone"`
	ID         int                 `json:"id"`
	Name       string              `json:"name"`
	Cod        int                 `json:"cod"`
}

type weatherMapCoord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type weatherMapWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type weatherMapMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type weatherMapWind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

type weatherMapClouds struct {
	All int `json:"all"`
}

type weatherMapSys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

// CurrentWeather returns the current weather for the given coordinate.
func (c *weatherMapClient) CurrentWeather(ctx context.Context, coordinate Coordinate) (*Weather, error) {
	req, err := c.newWeatherMapGetRequest(ctx, weatherV25Path)
	if err != nil {
		return nil, err
	}
	queryParams := req.URL.Query()
	queryParams.Set("lat", fmt.Sprintf("%f", coordinate.Lat))
	queryParams.Set("lon", fmt.Sprintf("%f", coordinate.Lon))
	req.URL.RawQuery = queryParams.Encode()

	weatherMapResponse := currentWeatherMapResponse{}
	err = c.executeRequest(req, &weatherMapResponse)
	if err != nil {
		return nil, err
	}

	description := ""
	if len(weatherMapResponse.Weather) > 0 {
		description = weatherMapResponse.Weather[0].Description
	}

	return &Weather{
		Description: description,
		WindSpeed:   weatherMapResponse.Wind.Speed,
		Temperature: weatherMapResponse.Main.Temp,
		Timestamp:   int64(weatherMapResponse.Dt),
	}, nil
}

// executeRequest executes the request and unmarshals the response into the response object.
func (c *weatherMapClient) executeRequest(req *http.Request, response any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Default().Println("error closing response body:", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedStatusCode{
			StatusCode: resp.StatusCode,
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return err
	}
	return nil
}

// newWeatherMapGetRequest creates a new http request with the api key inside the query params.
func (c *weatherMapClient) newWeatherMapGetRequest(ctx context.Context, apiPath string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s", path.Join(weatherMapBaseURL, apiPath)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	queryParams.Set("appid", c.apiKey)
	req.URL.RawQuery = queryParams.Encode()
	return req, nil
}
