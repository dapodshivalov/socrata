package cli

import (
	"fmt"
	"os"

	"socrata/src/socrata"
	"socrata/src/socrata/metrics"
	"socrata/src/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Render(topByMetricTypeAndYear map[metrics.MetricType]map[socrata.Year][]*metrics.Metric) {
	header := tableHeader(topByMetricTypeAndYear)
	printer := message.NewPrinter(language.English)
	for _, metricType := range metrics.CollectedMetrics {
		metricsByYear := topByMetricTypeAndYear[metricType]

		tableWriter := table.NewWriter()
		tableWriter.SetOutputMirror(os.Stdout)

		tableWriter.SetTitle(string(metricType))
		tableWriter.AppendHeader(header)

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

func tableHeader(topByMetricTypeAndYear map[metrics.MetricType]map[socrata.Year][]*metrics.Metric) table.Row {
	limit := 0
	for _, metricType := range metrics.CollectedMetrics {
		metricsByYear := topByMetricTypeAndYear[metricType]
		for _, metrics := range metricsByYear {
			limit = util.Max(limit, len(metrics))
		}
	}
	header := table.Row{"Year"}
	for i := range limit {
		header = append(header, fmt.Sprintf("#%d", i+1))
	}
	return header
}
