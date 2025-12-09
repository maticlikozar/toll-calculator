package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	targetURL = "http://toll.test/api/v1/toll-events"
	apiKey    = "123"
)

var vehicleTypes = []string{"car", "truck", "motorbike", "tractor"}

type TollEvent struct {
	LicensePlate string `json:"license_plate"`
	VehicleType  string `json:"vehicle_type"`
	EventStart   string `json:"event_start"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	client := &http.Client{Timeout: 10 * time.Second}

	// Command-line flags for max requests OR duration
	totalRequests := flag.Int("requests", 100000, "Total requests to send")
	rps := flag.Int("rps", 200, "Requests per second")
	durationSec := flag.Int("duration", 0, "Duration of test in seconds (overrides requests)")
	flag.Parse()

	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second / time.Duration(*rps))
	defer ticker.Stop()

	// Metrics
	var mu sync.Mutex
	var successCount, failCount int
	var totalLatency time.Duration
	var minLatency = time.Hour
	var maxLatency time.Duration

	fmt.Printf("Starting load test with RPS=%d…\n", *rps)

	startTime := time.Now()
	sent := 0
	var stopTime time.Time
	if *durationSec > 0 {
		stopTime = startTime.Add(time.Duration(*durationSec) * time.Second)
	}

	// range over ticker.C triggers on every tick
	for range ticker.C {
		// Exit conditions
		if *durationSec > 0 && time.Now().After(stopTime) {
			break
		}
		if *durationSec == 0 && sent >= *totalRequests {
			break
		}

		sent++
		wg.Add(1)

		go func() {
			defer wg.Done()

			event := TollEvent{
				LicensePlate: randomLicensePlate(),
				VehicleType:  randomVehicleType(),
				EventStart:   randomTimeToday().Format(time.RFC3339),
			}

			bodyBytes, _ := json.Marshal(event)
			req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(bodyBytes))
			if err != nil {
				mu.Lock()
				failCount++
				mu.Unlock()
				fmt.Println("request error:", err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-API-Key", apiKey)

			reqStart := time.Now()
			resp, err := client.Do(req)
			reqDuration := time.Since(reqStart)

			mu.Lock()
			totalLatency += reqDuration
			if reqDuration < minLatency {
				minLatency = reqDuration
			}
			if reqDuration > maxLatency {
				maxLatency = reqDuration
			}
			if err != nil {
				failCount++
				mu.Unlock()
				fmt.Println("http error:", err)
				return
			}
			if resp.StatusCode >= 300 {
				failCount++
			} else {
				successCount++
			}
			mu.Unlock()
			resp.Body.Close()
		}()
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	avgLatency := time.Duration(0)
	if successCount > 0 {
		avgLatency = totalLatency / time.Duration(successCount)
	}

	// Report
	fmt.Println("\n================== LOAD TEST REPORT ==================")
	fmt.Printf("Total requests sent       : %d\n", sent)
	fmt.Printf("Successful requests       : %d\n", successCount)
	fmt.Printf("Failed requests           : %d\n", failCount)
	fmt.Printf("Total test duration       : %v\n", totalDuration)
	fmt.Printf("Average request latency   : %v\n", avgLatency)
	fmt.Printf("Minimum request latency   : %v\n", minLatency)
	fmt.Printf("Maximum request latency   : %v\n", maxLatency)
	fmt.Printf("Achieved RPS             : %.2f\n", float64(sent)/totalDuration.Seconds())
	fmt.Println("=====================================================")
}

// Random license plate, e.g., "AB-123-XY"
func randomLicensePlate() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits := []rune("0123456789")

	return fmt.Sprintf("%c%c-%d%d%d-%c%c",
		letters[rand.Intn(len(letters))],
		letters[rand.Intn(len(letters))],
		digits[rand.Intn(len(digits))],
		digits[rand.Intn(len(digits))],
		digits[rand.Intn(len(digits))],
		letters[rand.Intn(len(letters))],
		letters[rand.Intn(len(letters))],
	)
}

func randomVehicleType() string {
	return vehicleTypes[rand.Intn(len(vehicleTypes))]
}

// Random time today, not in the future
func randomTimeToday() time.Time {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	secondsSinceStart := rand.Int63n(int64(now.Sub(startOfDay).Seconds()) + 1)
	return startOfDay.Add(time.Duration(secondsSinceStart) * time.Second)
}
