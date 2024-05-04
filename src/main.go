package main

import (
	"log"
	"os"

	"socrata/src/socrata/cli"
	"socrata/src/socrata/client"
	"socrata/src/socrata/metrics"

	"github.com/joho/godotenv"
)

const LIMIT = 10

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	token, ok := os.LookupEnv("SOCRATA_TOKEN")
	if !ok {
		log.Printf("No token provided. Using default requests quota.")
	}

	socrataClient := client.NewClient(token)
	realEstateTransactions, err := socrataClient.GetTransactions()
	if err != nil {
		log.Fatal(err)
	}

	topByMetricTypeAndYear := metrics.TopByMetricTypeAndYear(realEstateTransactions, LIMIT)

	cli.Render(topByMetricTypeAndYear)
}
