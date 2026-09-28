package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	lat, _ := strconv.ParseFloat(os.Getenv("LAT"), 64)
	lon, _ := strconv.ParseFloat(os.Getenv("LON"), 64)
	dist, _ := strconv.Atoi(os.Getenv("DIST"))

	url := fmt.Sprintf(
		"https://opendata.adsb.fi/api/v3/lat/%f/lon/%f/dist/%d",
		lat, lon, dist,
	)

	data, err := fetchAircraftData(url)
	if err != nil {
		fmt.Errorf(
			"Error while fetching aircraft data: %w ",
			err,
		)
	}

	for _, aircraft := range data {
		aircraft.Callsign = strings.TrimSpace(aircraft.Callsign)
		url := fmt.Sprintf("https://api.adsbdb.com/v0/callsign/%s", aircraft.Callsign)
		flightPlan, err := fetchFlightPlan(url)
		if err != nil {
			fmt.Errorf(
				"Error while fetching aircraft data: %w ",
				err,
			)
		}
	}

	fmt.Println("Aircraft data:", data)

}
