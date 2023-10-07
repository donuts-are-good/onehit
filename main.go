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
	// create a new key-value store
	store := libkeva.NewKeyValueStore("data.json", 5*time.Second)

	// load initial data from file if it exists
	if err := store.LoadFromFile("data.json"); err != nil {
		log.Fatalf("Failed to load data from file: %v", err)
	}

	// start a goroutine to print statistics
	go printStats(store)

	// create http server with cors middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		switch r.Method {
		case http.MethodGet:
			// check if path starts with "/x/"
			if len(r.URL.Path) > 3 && r.URL.Path[:3] == "/x/" {
				key := r.URL.Path[3:]
				// insert or update row with key
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

	// start http server
	if err := http.ListenAndServe(":3589", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func printStats(store *libkeva.KeyValueStore) {
	// start goroutine that prints statistics every 1 second
	for {
		data := store.GetData()
		fmt.Print("\033[2J\n" + time.Now().String()[:19] + "\n")
		fmt.Printf("Total keys: %d\n\n", len(data))

		// create tabwriter with padding of 4 spaces
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

		// print header row
		fmt.Fprintln(writer, "Key\tHits")

		// print rows for all keys
		for key, value := range data {
			fmt.Fprintf(writer, "%s\t%d\n", key, value)
		}

		// flush tabwriter buffer
		writer.Flush()

		time.Sleep(1 * time.Second)
	}
}
