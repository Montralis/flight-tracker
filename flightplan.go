package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type QueryResponse struct {
	Response FlightRoute `json:"response"`
}

type FlightRoute struct {
	FlightRoute FlightInfo `json:"flightroute"`
}

type FlightInfo struct {
	Callsign    string  `json:"callsign"`
	Airline     Airline `json:"airline"`
	Origin      Airport `json:"origin"`
	Destination Airport `json:"destination"`
}

type Airline struct {
	Name string `json:"name"`
}
type Airport struct {
	Name string `json:"name"`
	IATA string `json:"iata_code"`
	ICAO string `json:"icao_code"`
}

func fetchFlightPlan(url string) (*FlightInfo, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var query QueryResponse
	if err := json.Unmarshal(body, &query); err != nil {
		return nil, err
	}

	return &query.Response.FlightRoute, nil
}
