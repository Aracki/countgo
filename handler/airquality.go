package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const airPollutionAPIFormat = "http://api.openweathermap.org/data/2.5/air_pollution?lat=44.7866&lon=20.4489&appid=%s"

var belgradeTZ *time.Location

func init() {
	var err error
	belgradeTZ, err = time.LoadLocation("Europe/Belgrade")
	if err != nil {
		belgradeTZ = time.UTC
	}
}

type airPollutionResponse struct {
	List []struct {
		Main struct {
			AQI int `json:"aqi"`
		} `json:"main"`
		Components struct {
			PM25 float64 `json:"pm2_5"`
		} `json:"components"`
	} `json:"list"`
}

type airQualityCache struct {
	sync.RWMutex
	aqi       int
	pm25      float64
	updatedAt time.Time
	valid     bool
}

var aqCache = &airQualityCache{}

func StartAirQualityUpdater() {
	// Fetch immediately on startup
	fetchAirQuality()

	// Then fetch every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			fetchAirQuality()
		}
	}()
}

func fetchAirQuality() {
	apiKey := viper.GetString("openweathermap.apikey")
	if apiKey == "" {
		log.Println("fetchAirQuality: openweathermap.apikey not configured")
		return
	}
	resp, err := http.Get(fmt.Sprintf(airPollutionAPIFormat, apiKey))
	if err != nil {
		log.Println("fetchAirQuality: api error: ", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("fetchAirQuality: read error: ", err.Error())
		return
	}

	var data airPollutionResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("fetchAirQuality: json error: ", err.Error())
		return
	}

	if len(data.List) == 0 {
		log.Println("fetchAirQuality: no data in response")
		return
	}

	aqCache.Lock()
	aqCache.aqi = data.List[0].Main.AQI
	aqCache.pm25 = data.List[0].Components.PM25
	aqCache.updatedAt = time.Now()
	aqCache.valid = true
	aqCache.Unlock()

	log.Printf("fetchAirQuality: updated AQI=%d PM2.5=%.1f", aqCache.aqi, aqCache.pm25)
}

func aqiLabel(aqi int) string {
	switch aqi {
	case 1:
		return "Good"
	case 2:
		return "Fair"
	case 3:
		return "Moderate"
	case 4:
		return "Poor"
	case 5:
		return "Very Poor"
	default:
		return "Unknown"
	}
}

func aqiColor(aqi int) string {
	switch aqi {
	case 1:
		return "#4CAF50"
	case 2:
		return "#8BC34A"
	case 3:
		return "#FFC107"
	case 4:
		return "#FF9800"
	case 5:
		return "#F44336"
	default:
		return "#9E9E9E"
	}
}

func airQuality(w http.ResponseWriter, r *http.Request) {
	aqCache.RLock()
	if !aqCache.valid {
		aqCache.RUnlock()
		http.Error(w, "Air quality data not available yet", http.StatusServiceUnavailable)
		return
	}
	aqi := aqCache.aqi
	pm25 := aqCache.pm25
	updatedAt := aqCache.updatedAt
	aqCache.RUnlock()

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Belgrade Air Quality</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%%, #16213e 100%%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            color: #fff;
        }
        .container {
            text-align: center;
            padding: 2rem;
        }
        h1 {
            font-size: 1.5rem;
            font-weight: 300;
            margin-bottom: 2rem;
            opacity: 0.9;
        }
        .cards {
            display: flex;
            gap: 2rem;
            flex-wrap: wrap;
            justify-content: center;
        }
        .card {
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 2rem 3rem;
            min-width: 200px;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        .card-label {
            font-size: 0.875rem;
            text-transform: uppercase;
            letter-spacing: 2px;
            opacity: 0.7;
            margin-bottom: 0.5rem;
        }
        .card-value {
            font-size: 4rem;
            font-weight: 700;
            line-height: 1;
        }
        .card-unit {
            font-size: 0.875rem;
            opacity: 0.7;
            margin-top: 0.5rem;
        }
        .aqi-status {
            font-size: 0.875rem;
            margin-top: 1rem;
            opacity: 0.7;
        }
        .location {
            margin-top: 2rem;
            opacity: 0.5;
            font-size: 0.875rem;
        }
        .updated {
            margin-top: 0.5rem;
            opacity: 0.5;
            font-size: 0.875rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Belgrade Air Quality</h1>
        <div class="cards">
            <div class="card">
                <div class="card-label">Air Quality Index</div>
                <div class="card-value" style="color: %s">%d</div>
                <div class="aqi-status">%s</div>
            </div>
            <div class="card">
                <div class="card-label">PM2.5</div>
                <div class="card-value">%.1f</div>
                <div class="card-unit">&micro;g/m&sup3;</div>
            </div>
        </div>
        <div class="location">44.7866&deg;N, 20.4489&deg;E</div>
        <div class="updated">Last updated at %s</div>
    </div>
</body>
</html>`, aqiColor(aqi), aqi, aqiLabel(aqi), pm25, updatedAt.In(belgradeTZ).Format("15:04:05"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
