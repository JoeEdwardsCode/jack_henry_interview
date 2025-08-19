package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client for National Weather Service
// http://weather.gov/documentation/services-web-api
const NWSHost = "https://api.weather.gov/"

type NWSClient struct {
	BaseURL     string
	Client      *http.Client
	ColdMax     int
	ModerateMax int
}

func NewNWSClient(conf Config) *NWSClient {
	return &NWSClient{
		BaseURL: conf.NWSBaseURL,
		Client: &http.Client{
			Timeout: conf.RequestTimeout,
		},
		ColdMax:     conf.ColdMax,
		ModerateMax: conf.ModerateMax,
	}
}

func (c *NWSClient) Get(endpoint string) ([]byte, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Request creation error: %w", err)
	}

	req.Header.Set("User-Agent", "nws_client")
	req.Header.Set("Accept", "application/geo+json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Response body read error: %w", err)
	}

	return body, nil
}

// GetWeatherData calls various NWS endpoints to gerate temperature summary from coordinates
func (c *NWSClient) GetWeatherData(r WeatherRequest) WeatherResponse {
	endpoint := fmt.Sprintf("points/%s,%s", r.Location.Latitude, r.Location.Longitude)
	pointBytes, err := c.Get(endpoint)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to reach NWS points API: %v", err)
		return WeatherResponse{
			Status:       WeatherResponseStatusFailure,
			ErrorMessage: &errorMessage,
		}
	}

	var pointResponse NWSPointResponse
	err = json.Unmarshal(pointBytes, &pointResponse)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to parse NWS point response: %v", err)
		return WeatherResponse{
			Status:       WeatherResponseStatusFailure,
			ErrorMessage: &errorMessage,
		}
	}

	forecastURL := strings.TrimPrefix(pointResponse.Properties.Forecast, "https://api.weather.gov/")
	forecastBytes, err := c.Get(forecastURL)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to reach NWS forecast API: %v", err)
		return WeatherResponse{
			Status:       WeatherResponseStatusFailure,
			ErrorMessage: &errorMessage,
		}
	}

	var forecastResponse NWSForecastResponse
	err = json.Unmarshal(forecastBytes, &forecastResponse)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to parse NWS forecast response: %v", err)
		return WeatherResponse{
			Status:       WeatherResponseStatusFailure,
			ErrorMessage: &errorMessage,
		}
	}

	if len(forecastResponse.Properties.Periods) == 0 {
		errorMessage := "No forecast periods available"
		return WeatherResponse{
			Status:       WeatherResponseStatusFailure,
			ErrorMessage: &errorMessage,
		}
	}

	currentPeriod := forecastResponse.Properties.Periods[0]

	var characterization string
	if currentPeriod.Temperature < c.ColdMax {
		characterization = "Cold"
	} else if currentPeriod.Temperature < c.ModerateMax {
		characterization = "Moderate"
	} else {
		characterization = "Hot"
	}

	return WeatherResponse{
		Status:           WeatherResponseStatusSuccess,
		Forecast:         &currentPeriod.ShortForecast,
		Characterization: &characterization,
	}
}
