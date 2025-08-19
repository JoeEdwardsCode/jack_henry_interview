package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherServer struct {
	client *NWSClient
	config *Config
}

func NewWeatherServer(config *Config) *WeatherServer {
	client := NewNWSClient(*config)

	return &WeatherServer{
		client: client,
		config: config,
	}
}

func (s *WeatherServer) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, fmt.Sprintf("%s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req WeatherRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		s.writeError(w, fmt.Sprintf("%s", err.Error()), http.StatusBadRequest)
		return
	}

	response := s.client.GetWeatherData(req)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *WeatherServer) writeError(w http.ResponseWriter, errorCode string, statusCode int) {
	response := WeatherResponse{
		Status:       "failure",
		ErrorMessage: &errorCode,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (s *WeatherServer) Start() error {
	http.HandleFunc("/weather", s.WeatherHandler)

	fmt.Printf("server listening on port %s\n", s.config.ServerPort)
	fmt.Printf("Cold threshold: %d°F, Moderate threshold: %d°F\n", s.config.ColdMax, s.config.ModerateMax)

	return http.ListenAndServe(":"+s.config.ServerPort, nil)
}
