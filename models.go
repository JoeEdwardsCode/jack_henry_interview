package main

const (
	WeatherResponseStatusFailure = "failure"
	WeatherResponseStatusSuccess = "success"
)

type Location struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type WeatherRequest struct {
	Location Location `json:"location"`
}

type WeatherResponse struct {
	Status           string  `json:"status"`
	Forecast         *string `json:"forecast,omitempty"`
	Characterization *string `json:"characterization,omitempty"`
	ErrorMessage     *string `json:"error_code,omitempty"`
}

// NWS API response structures
type NWSPointResponse struct {
	Properties NWSPointProperties `json:"properties"`
}

type NWSPointProperties struct {
	Forecast         string              `json:"forecast"`
	ForecastHourly   string              `json:"forecastHourly"`
	GridID           string              `json:"gridId"`
	GridX            int                 `json:"gridX"`
	GridY            int                 `json:"gridY"`
	ForecastOffice   string              `json:"forecastOffice"`
	RelativeLocation NWSRelativeLocation `json:"relativeLocation"`
}

type NWSRelativeLocation struct {
	Properties NWSLocationProperties `json:"properties"`
}

type NWSLocationProperties struct {
	City  string `json:"city"`
	State string `json:"state"`
}

type NWSForecastResponse struct {
	Properties NWSForecastProperties `json:"properties"`
}

type NWSForecastProperties struct {
	Periods []NWSForecastPeriod `json:"periods"`
}

type NWSForecastPeriod struct {
	Number           int    `json:"number"`
	Name             string `json:"name"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	IsDaytime        bool   `json:"isDaytime"`
	Temperature      int    `json:"temperature"`
	TemperatureUnit  string `json:"temperatureUnit"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}
