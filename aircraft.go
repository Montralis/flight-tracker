package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type AirCraftQueryResponse struct {
	Aircrafts []Aircraft `json:"ac"`
}
type Aircraft struct {
	Callsign string  `json:"flight"`
	Speed    float64 `json:"gs"`
	Type     string  `json:"desc"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Dist     float64 `json:"dst"`
}

func fetchAircraftData(url string) ([]Aircraft, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var query AirCraftQueryResponse
	if err := json.Unmarshal(body, &query); err != nil {
		return nil, err
	}

	return query.Aircrafts, nil
}
