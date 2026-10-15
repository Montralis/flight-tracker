package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"golang.org/x/exp/slog"
)

type Track struct {
	Callsign string
	Aircraft string
	From     string
	To       string
}

func flightsInSquare(writer http.ResponseWriter, request *http.Request) {
	ul, _ := strconv.ParseFloat(request.PathValue("ul"), 64)
	ur, _ := strconv.ParseFloat(request.PathValue("ur"), 64)
	ll, _ := strconv.ParseFloat(request.PathValue("ll"), 64)
	lr, _ := strconv.ParseFloat(request.PathValue("lr"), 64)
	cLat := (ul + ll) / 2
	cLon := (ur + lr) / 2

	latDist := (ul - ll) * 60
	lonDist := (ur - lr) * 60 * math.Cos(cLat*math.Pi/180)

	dist := math.Sqrt(
		math.Pow(latDist/2, 2) +
			math.Pow(lonDist/2, 2),
	)

	polygon := orb.Polygon{
		{
			{lr, ll},
			{ur, ll},
			{ur, ul},
			{lr, ul},
			{lr, ll},
		},
	}

	url := fmt.Sprintf(
		"https://opendata.adsb.fi/api/v3/lat/%f/lon/%f/dist/%f",
		cLat, cLon, dist,
	)

	var TrackStore []Track

	aircrafts, err := fetchAircraftData(url)
	if err != nil {
		slog.Error("Error while fetching aircraft data", "err", err)
		http.Error(writer, "Failed to fetch aircraft", http.StatusInternalServerError)
		return
	}

	for _, aircraft := range aircrafts {

		point := orb.Point{
			aircraft.Lon,
			aircraft.Lat,
		}

		isIn := planar.PolygonContains(polygon, point)

		if !isIn {
			continue
		}

		aircraft.Callsign = strings.TrimSpace(aircraft.Callsign)
		url := fmt.Sprintf("https://api.adsbdb.com/v0/callsign/%s", aircraft.Callsign)
		flightPlan, err := fetchFlightPlan(url)
		if err != nil {
			slog.Error("Error while fetching flightplan data", "err", err)
			http.Error(writer, "Failed to fetch flightplan", http.StatusInternalServerError)
			return
		}
		if flightPlan != nil {
			TrackStore = append(TrackStore, Track{
				Aircraft: aircraft.Type,
				Callsign: aircraft.Callsign,
				From:     flightPlan.Origin.IATA,
				To:       flightPlan.Destination.IATA,
			})
		}
	}

	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(TrackStore); err != nil {
		slog.Error(
			"failed to encode flight overview",
			slog.String("err", err.Error()),
		)
	}
}

func flightsInCircle(writer http.ResponseWriter, request *http.Request) {
	lat := request.PathValue("lat")
	lon := request.PathValue("lon")
	dist := request.PathValue("dist")

	url := fmt.Sprintf(
		"https://opendata.adsb.fi/api/v3/lat/%s/lon/%s/dist/%s",
		lat, lon, dist,
	)

	var TrackStore []Track

	aircrafts, err := fetchAircraftData(url)
	if err != nil {
		slog.Error("Error while fetching aircraft data", "err", err)
		http.Error(writer, "Failed to fetch aircraft", http.StatusInternalServerError)
		return
	}

	for _, aircraft := range aircrafts {
		aircraft.Callsign = strings.TrimSpace(aircraft.Callsign)
		url := fmt.Sprintf("https://api.adsbdb.com/v0/callsign/%s", aircraft.Callsign)

		flightPlan, err := fetchFlightPlan(url)

		if err != nil {
			slog.Error("Error while fetching flightplan data", "err", err)
			http.Error(writer, "Failed to fetch flightplan", http.StatusInternalServerError)
			return
		}

		if flightPlan != nil {
			TrackStore = append(TrackStore, Track{
				Aircraft: aircraft.Type,
				Callsign: aircraft.Callsign,
				From:     flightPlan.Origin.IATA,
				To:       flightPlan.Destination.IATA,
			})
		}
	}

	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(TrackStore); err != nil {
		slog.Error(
			"failed to encode flight overview",
			slog.String("err", err.Error()),
		)
	}
}

func main() {

	_ = godotenv.Load()
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/flights/{lat}/{lon}/{dist}", flightsInCircle)
	mux.HandleFunc("GET /api/flights/{ul}/{ur}/{ll}/{lr}", flightsInSquare)

	listenAddress := ":" + strconv.Itoa(serverPort)

	httpServer := &http.Server{
		Addr:    listenAddress,
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"listen failed",
			slog.String("err", err.Error()),
		)
		os.Exit(1)
	}
}
