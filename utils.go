package main

import (
	"fmt"
	"strconv"
	"strings"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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

func getStats(db *sqlx.DB) (keys int, hits int, popularKeys []KeyValue) {
	err := db.Get(&keys, "SELECT COUNT(*) FROM kv")
	if err != nil {
		log.Println(err)
	}

	err = db.Get(&hits, "SELECT SUM(value) FROM kv")
	if err != nil {
		log.Println(err)
	}

	err = db.Select(&popularKeys, "SELECT * FROM kv ORDER BY value DESC LIMIT 10")
	if err != nil {
		log.Println(err)
	}

	return
}
