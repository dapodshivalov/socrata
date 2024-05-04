package main

import (
	"fmt"
	"log"
	"os"

	"socrata/src/socrata"
	"socrata/src/socrata/client"
	"socrata/src/socrata/metrics"
	"socrata/src/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joho/godotenv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const LIMIT = 10

var topHeader table.Row

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	topHeader = make([]interface{}, LIMIT)
	for i := range LIMIT {
		topHeader[i] = fmt.Sprintf("#%d", i)
	}
	topHeader = append(table.Row{"Year"}, topHeader...)
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

	renderResult(topByMetricTypeAndYear)
}

func renderResult(topByMetricTypeAndYear map[metrics.MetricType]map[socrata.Year][]*metrics.Metric) {
	printer := message.NewPrinter(language.English)
	for _, metricType := range metrics.CollectedMetrics {
		metricsByYear := topByMetricTypeAndYear[metricType]
		tableWriter := table.NewWriter()
		tableWriter.SetOutputMirror(os.Stdout)
		tableWriter.SetTitle(string(metricType))
		tableWriter.AppendHeader(topHeader)
		years := util.SortedMapKeys(metricsByYear, func(a, b socrata.Year) bool {
			return a > b
		})
		rows := make([]table.Row, 0, len(years))
		for _, year := range years {
			row := table.Row{year}
			for _, townAndMetricsValue := range metricsByYear[year] {
				row = append(row, printer.Sprintf("%s\n%.2f", townAndMetricsValue.Town, townAndMetricsValue.Value))
			}
			rows = append(rows, row)
		}
		tableWriter.AppendRows(rows)
		tableWriter.Render()
	}
}
