package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr, ok := os.LookupEnv("CONNECTION_STRING")
	if !ok {
		log.Fatal("missing CONNECTION_STRING env var")
	}

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", handleGetRegion(db))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetRegion(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRow(context.Background(), "SELECT gateway_region()")

		var region string
		if err := row.Scan(&region); err != nil {
			fmt.Fprintf(w, "error querying gateway_region: %v\n", err)
		}

		fmt.Fprintln(w, region)
	}
}
