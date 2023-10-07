package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/donuts-are-good/libkeva"
)

func main() {
	// Create a new KeyValueStore
	store := libkeva.NewKeyValueStore("data.json", 5*time.Second)

	// Load initial data from file if it exists
	err := store.LoadFromFile("data.json")
	if err != nil {
		log.Fatal(err)
	}

	go printStats(store)

	// Create HTTP server with CORS middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		switch r.Method {
		case http.MethodGet:
			// Check if path starts with "/x/"
			if len(r.URL.Path) > 3 && r.URL.Path[:3] == "/x/" {
				key := r.URL.Path[3:]
				// Insert or update row with key
				value, exists := store.Get(key)
				if !exists {
					value = 0
				}
				store.Set(key, value.(int)+1)
				fmt.Fprintf(w, "Registered hit for key %s", key)
			} else {
				http.NotFound(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":3589", handler))
}

func printStats(store *libkeva.KeyValueStore) {
	// Start goroutine that prints statistics every 1 second
	for {
		data := store.GetData()
		fmt.Print("\033[2J\n" + time.Now().String()[:19] + "\n")
		fmt.Printf("Total keys: %d\n\n", len(data))

		// Create tabwriter with padding of 4 spaces
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

		// Print header row
		fmt.Fprintln(writer, "Key\tHits")

		// Print rows for all keys
		for key, value := range data {
			fmt.Fprintf(writer, "%s\t%d\n", key, value)
		}

		// Flush tabwriter buffer
		writer.Flush()

		time.Sleep(1 * time.Second)
	}
}
