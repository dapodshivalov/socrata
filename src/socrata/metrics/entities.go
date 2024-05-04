package metrics

import (
	"socrata/src/socrata"
	"socrata/src/util"
)

type MetricType string

const (
	AVG_SALES_RATIO MetricType = "avg_sales_ratio"
	AVG_VOLUME      MetricType = "avg_volume"
)

var CollectedMetrics = []MetricType{
	AVG_SALES_RATIO,
	AVG_VOLUME,
}

// ----------------------------------------

type Metric struct {
	Type  MetricType
	Town  socrata.Town
	Value float64
}

func MakeMetric(t MetricType, town socrata.Town, value float64) *Metric {
	return &Metric{
		Type:  t,
		Town:  town,
		Value: value,
	}
}

// ----------------------------------------

type MetricsMap map[MetricType]*Metric

func MakeMetricsMap(metrics ...*Metric) MetricsMap {
	return util.SliceToMap(metrics, func(metric *Metric) MetricType {
		return metric.Type
	})
}
