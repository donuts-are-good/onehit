package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    "github.com/jmoiron/sqlx"
		_ "github.com/mattn/go-sqlite3"

)

type KeyValue struct {
    Key   string `db:"key"`
    Value int    `db:"value"`
}

func main() {
    // Open SQLite database
    db, err := sqlx.Connect("sqlite3", "kv.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create table if it doesn't exist
    db.MustExec("CREATE TABLE IF NOT EXISTS kv (key TEXT PRIMARY KEY, value INTEGER)")

    // Create HTTP server with CORS middleware
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")

        switch r.Method {
        case http.MethodGet:
            // Check if path starts with "/x/"
            if len(r.URL.Path) > 3 && r.URL.Path[:3] == "/x/" {
                key := r.URL.Path[3:]
                // Insert or update row with key
                db.MustExec("INSERT OR REPLACE INTO kv (key, value) VALUES (?, COALESCE((SELECT value FROM kv WHERE key = ?), 0) + 1)", key, key)
                fmt.Fprintf(w, "Registered hit for key %s", key)
            } else if len(r.URL.Path) > 3 && r.URL.Path[:3] == "/r/" {
                key := r.URL.Path[3:]
                // Get row with key
                var kv KeyValue
                err := db.Get(&kv, "SELECT * FROM kv WHERE key = ?", key)
                if err == sql.ErrNoRows {
                    fmt.Fprintf(w, "Key %s not found", key)
                } else if err != nil {
                    log.Println(err)
                    http.Error(w, "Internal server error", http.StatusInternalServerError)
                } else {
                    // Format value as human-readable string
                    valueStr := formatValue(kv.Value)
                    fmt.Fprint(w, valueStr)
                }
            } else {
                http.NotFound(w, r)
            }
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Start HTTP server
    log.Fatal(http.ListenAndServe(":8080", handler))
}

func formatValue(value int) string {
    var suffix string
    var formattedValue string

    if value >= 1000000000 {
        suffix = "B"
        formattedValue = fmt.Sprintf("%.2f", float64(value)/1000000000)
    } else if value >= 1000000 {
        suffix = "M"
        formattedValue = fmt.Sprintf("%.2f", float64(value)/1000000)
    } else if value >= 1000 {
        suffix = "k"
        formattedValue = fmt.Sprintf("%.1f", float64(value)/1000)
    } else {
        return strconv.Itoa(value)
    }

    // Remove trailing ".00" from formatted value
    formattedValue = strings.TrimRight(formattedValue, "0")

    // Remove trailing "." if present
    formattedValue = strings.TrimSuffix(formattedValue, ".")

    return formattedValue + suffix
}
