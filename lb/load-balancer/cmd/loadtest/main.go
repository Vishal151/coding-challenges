package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "http://localhost:8080", "URL to test")
	concurrency := flag.Int("c", 10, "Number of concurrent requests")
	totalRequests := flag.Int("n", 1000, "Total number of requests")
	flag.Parse()

	var wg sync.WaitGroup
	results := make(chan float64, *totalRequests)

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < *totalRequests / *concurrency; j++ {
				start := time.Now()
				resp, err := http.Get(*url)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}
				resp.Body.Close()
				duration := time.Since(start).Seconds()
				results <- duration
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalDuration float64
	var count int
	for duration := range results {
		totalDuration += duration
		count++
	}

	avgDuration := totalDuration / float64(count)
	fmt.Printf("Average request duration: %.4f seconds\n", avgDuration)
	fmt.Printf("Requests per second: %.2f\n", float64(count)/totalDuration)
}
