package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type KeyValue struct {
	Key    string `db:"key"`
	Values []byte `db:"values"`
}

const (
	defaultWindowSize = 30
	maxWindowSize     = 365
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := http.NewServeMux()
	router.HandleFunc("/x/", handleRegisterHit(db))
	router.HandleFunc("/r/", handleGetHitCount(db))

	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(router)))
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func initDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "kv.db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.MustExec("CREATE TABLE IF NOT EXISTS kv (key TEXT PRIMARY KEY, hit_counts BLOB)")
	return db, nil
}


func handleRegisterHit(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if len(r.URL.Path) <= 3 || r.URL.Path[:3] != "/x/" {
			http.NotFound(w, r)
			return
		}

		key := r.URL.Path[3:]
		now := time.Now()
		day := now.Day()
		value := []byte{0}
		var kv KeyValue
		err := db.Get(&kv, "SELECT * FROM kv WHERE key = ?", key)
		if err == sql.ErrNoRows {
			buf := new(bytes.Buffer)
			enc := gob.NewEncoder(buf)
			enc.Encode(value)
			db.MustExec("INSERT INTO kv (key, values) VALUES (?, ?)", key, buf.Bytes())
		} else if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			buf := bytes.NewBuffer(kv.Values)
			dec := gob.NewDecoder(buf)
			dec.Decode(&value)
			value[day-1]++
			buf.Reset()
			enc := gob.NewEncoder(buf)
			enc.Encode(value)
			db.MustExec("UPDATE kv SET values = substr(values, 2) || ? WHERE key = ?", buf.Bytes(), key)
		}
		fmt.Fprintf(w, "Registered hit for key %s", key)
	}
}

func handleGetHitCount(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if len(r.URL.Path) <= 3 || r.URL.Path[:3] != "/r/" {
			http.NotFound(w, r)
			return
		}

		key := r.URL.Path[3:]
		var windowSize int
		if len(r.URL.Path) > 4 {
			size, err := strconv.Atoi(r.URL.Path[4:])
			if err == nil && size > 0 && size <= maxWindowSize {
				windowSize = size
			} else {
				windowSize = defaultWindowSize
			}
		} else {
			windowSize = defaultWindowSize
		}
		now := time.Now()
		startDay := now.AddDate(0, 0, -windowSize+1).Day()
		var kv KeyValue
		err := db.Get(&kv, "SELECT * FROM kv WHERE key = ?", key)
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "Key %s not found", key)
			return
		} else if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		buf := bytes.NewBuffer(kv.Values)
		dec := gob.NewDecoder(buf)
		var values []byte
		err = dec.Decode(&values)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		startIndex := startDay - 1
		if len(values) < startIndex {
			startIndex = len(values)
		}
		sum := 0
		for i := startIndex; i < len(values); i++ {
			sum += int(values[i])
		}
		fmt.Fprintf(w, "%d", sum)
	}
}
