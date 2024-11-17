package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
)

var (
	baseURL           = "http://localhost:7111" // Change as needed
	concurrentClients = 100                     // Default concurrency
)
var uuids []string

// Helper function to insert a document
func insertDocument(b *testing.B) error {
	doc, _, err := generateRandomDoc()
	if err != nil {
		return err
	}
	resp, err := http.Post(baseURL+"/insert", "application/json", bytes.NewReader(doc))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("ERROR: %+v\n", resp)
		return errors.New("wrong status code")
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	//     fmt.Printf("response %+v\n, data %v", resp, string(data))
	return err
}

// Helper function to search for a document by UUID
func searchDocument(b *testing.B, uuid string) error {
	resp, err := http.Get(fmt.Sprintf("%s/search?id=%s", baseURL, uuid))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

// Benchmark for insert endpoint
func BenchmarkDBInsert(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan error, b.N)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- insertDocument(b)
		}()
	}

	wg.Wait()
	close(ch)

	// Count errors
	errorCount := 0
	for err := range ch {
		if err != nil {
			errorCount++
		}
	}

	b.ReportMetric(float64(b.N-errorCount)/float64(b.N), "success_rate")
}

// Benchmark for search endpoint
func BenchmarkInsert(b *testing.B) {
	// Generate random documents and insert them
	errorCount := 0
	for i := 0; i < b.N; i++ {
		doc, uuid, err := generateRandomDoc()
		if err != nil {
			b.Fatal(err)
		}
		_, err = http.Post(baseURL+"/insert", "application/json", bytes.NewReader(doc))
		if err != nil {
			//             b.Fatal(err)
			errorCount++
		}
		uuids = append(uuids, uuid)
	}
	b.ReportMetric(float64(b.N-errorCount)/float64(b.N), "success_rate")
}

// Benchmark for search endpoint
func BenchmarkSearch(b *testing.B) {
	// Benchmark the search by UUID
	var wg sync.WaitGroup
	ch := make(chan error, b.N)

	if len(uuids) == 0 {
		err := loadConfig("config.json")
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		err = connectMongoDB()
		if err != nil {
			b.Fatal(err)
		}

		uuids, err = getAllIDs(config.MongoCollection, "id")
		if err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(uuid string) {
			defer wg.Done()
			ch <- searchDocument(b, uuid)
		}(uuids[i%len(uuids)]) // Reuse UUIDs if we have fewer than b.N

	}
	wg.Wait()
	close(ch)

	// Count errors
	errorCount := 0
	for err := range ch {
		if err != nil {
			errorCount++
		}
	}

	b.ReportMetric(float64(b.N-errorCount)/float64(b.N), "success_rate")
}
