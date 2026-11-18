# Flighttracker

A simple flight tracking API written in Go. The application retrieves live aircraft data
from [ADS-B.fi](https://adsb.fi/) and enriches it with flight route information from [ADSBDB](https://www.adsbdb.com/).

## Features

- Retrieve live aircraft data using ADS-B.fi.
- Search for aircraft within a geographic circle.
- Search for aircraft within a rectangular area.
- Retrieve additional flight information:
    - Callsign
    - Aircraft type
    - Origin airport
    - Destination airport
- JSON REST API.

## Requirements

- Go
- Internet connection

## Configuration

Create a `.env` file in the project root with:

```bash
cp .env.dist .env
```

## Run

Start the application direct from the terminal:

```bash
go run ./src
```

The API will be available at (read port from .env):

```text
http://localhost:8080
```

## Compile

**Linux ARM64:**

```bash
GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o Flight_Tracker_ARM64 ./src
```

**Linux AMD64 (x64):**

```bash
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o Flight_Tracker_AMD64 ./src
```

## API Endpoints

### Flights in a Circle

Retrieve aircraft within a specified radius around a geographic position.

```http
GET /api/flights/{lat}/{lon}/{dist}
```

Example:

```http
GET /api/flights/49.8728/8.6512/100
```

Parameters:

- `lat`: Latitude of the center point.
- `lon`: Longitude of the center point.
- `dist`: Search radius in nautical miles.

### Flights in a Rectangle

Retrieve aircraft within a rectangular geographic area.

```http
GET /api/flights/{ul}/{ur}/{ll}/{lr}
```

Parameters:

- `ul`: Upper-left latitude.
- `ur`: Upper-right longitude.
- `ll`: Lower-left latitude.
- `lr`: Lower-right longitude.

## Example Response

```json
[
  {
    "Callsign": "SIA325",
    "Aircraft": "BOEING 777-300ER",
    "From": "FRA",
    "To": "SIN"
  }
]
```

## Data Sources

- [ADS-B.fi](https://adsb.fi/) – Live aircraft positions and information.
- [ADSBDB](https://www.adsbdb.com/) – Flight route and airport information.

