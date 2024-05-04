package metrics

import (
	"sort"

	"socrata/src/socrata"
	"socrata/src/util"
)

func TopByMetricTypeAndYear(realEstateTransactions []*socrata.RealEstateTransaction, limit int) map[MetricType]map[socrata.Year][]*Metric {
	aggregatedMetricsByYearAndTown := make(map[util.Pair[socrata.Year, socrata.Town]]MetricsMap)
	groupByYearAndTown := make(map[util.Pair[socrata.Year, socrata.Town]][]*socrata.RealEstateTransaction)
	for _, transaction := range realEstateTransactions {
		yearAndTown := util.MakePair(transaction.SaleYear, transaction.Town)
		groupByYearAndTown[yearAndTown] = append(groupByYearAndTown[yearAndTown], transaction)
	}
	for yearAndTown, groupedTransactions := range groupByYearAndTown {
		town := yearAndTown.Second
		aggregatedMetricsByYearAndTown[yearAndTown] = MakeMetricsMap(
			MakeMetric(AVG_VOLUME, town, Average(groupedTransactions, (*socrata.RealEstateTransaction).GetSalesVolume)),
			MakeMetric(AVG_SALES_RATIO, town, Average(groupedTransactions, (*socrata.RealEstateTransaction).GetSalesRatio)),
		)
	}

	topByMetricTypeAndYear := make(map[MetricType]map[socrata.Year][]*Metric)

	for yearAndTown, aggregatedMetrics := range aggregatedMetricsByYearAndTown {
		for metricType, metric := range aggregatedMetrics {
			if _, ok := topByMetricTypeAndYear[metricType]; !ok {
				topByMetricTypeAndYear[metricType] = make(map[socrata.Year][]*Metric)
			}
			byYear := topByMetricTypeAndYear[metricType]

			year := yearAndTown.First
			if _, ok := byYear[year]; !ok {
				byYear[year] = make([]*Metric, 0)
			}
			byYear[year] = append(byYear[year], metric)
		}
	}

	for _, topByYears := range topByMetricTypeAndYear {
		for year := range topByYears {
			metrics := topByYears[year]
			sort.SliceStable(metrics, func(i, j int) bool {
				return metrics[i].Value > metrics[j].Value
			})
			topByYears[year] = metrics[:util.Min(limit, len(metrics))]
		}
	}
	return topByMetricTypeAndYear
}
