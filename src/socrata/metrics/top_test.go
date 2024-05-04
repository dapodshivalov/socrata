package metrics

import (
	"testing"

	"socrata/src/socrata"

	"github.com/stretchr/testify/assert"
)

const (
	IZH = socrata.Town("Izhevsk")
	MSC = socrata.Town("Moscow")
	LSB = socrata.Town("Lisbon")
	MMI = socrata.Town("Miami")
)

func TestSuccess(t *testing.T) {

	// arrange
	transactions := []*socrata.RealEstateTransaction{
		socrata.NewTransaction(MSC, 2000, 0.0, 1000.0),
		socrata.NewTransaction(MSC, 2009, 10.0, 10000.0),
		socrata.NewTransaction(MSC, 2009, 11.0, 10000.0),
		socrata.NewTransaction(IZH, 2000, 0.5, 1000.0),
		socrata.NewTransaction(IZH, 2000, 0.6, 1200.0),
		socrata.NewTransaction(IZH, 2001, 0.8, 300.0),
	}

	expected := map[MetricType]map[socrata.Year][]*Metric{
		AVG_VOLUME: {
			2000: {
				MakeMetric(AVG_VOLUME, IZH, 1100.0),
				MakeMetric(AVG_VOLUME, MSC, 1000.0),
			},
			2001: {
				MakeMetric(AVG_VOLUME, IZH, 300.0),
			},
			2009: {
				MakeMetric(AVG_VOLUME, MSC, 10000.0),
			},
		},
		AVG_SALES_RATIO: {
			2000: {
				MakeMetric(AVG_SALES_RATIO, IZH, 0.55),
				MakeMetric(AVG_SALES_RATIO, MSC, 0.0),
			},
			2001: {
				MakeMetric(AVG_SALES_RATIO, IZH, 0.8),
			},
			2009: {
				MakeMetric(AVG_SALES_RATIO, MSC, 10.5),
			},
		},
	}

	// act
	top := TopByMetricTypeAndYear(transactions, 10)

	// assert
	assert.Equal(t, expected, top)
}

func TestLimit(t *testing.T) {

	// arrange
	transactions := []*socrata.RealEstateTransaction{
		socrata.NewTransaction(MSC, 2000, 9.0, 3000.0),
		socrata.NewTransaction(IZH, 2000, 4.0, 1000.0),
		socrata.NewTransaction(LSB, 2000, 7.0, 7000.0),
		socrata.NewTransaction(MMI, 2000, 2.0, 16000.0),
	}

	expected := map[MetricType]map[socrata.Year][]*Metric{
		AVG_VOLUME: {
			2000: {
				MakeMetric(AVG_VOLUME, MMI, 16000.0),
				MakeMetric(AVG_VOLUME, LSB, 7000.0),
				MakeMetric(AVG_VOLUME, MSC, 3000.0),
			},
		},
		AVG_SALES_RATIO: {
			2000: {
				MakeMetric(AVG_SALES_RATIO, MSC, 9.0),
				MakeMetric(AVG_SALES_RATIO, LSB, 7.0),
				MakeMetric(AVG_SALES_RATIO, IZH, 4.0),
			},
		},
	}

	// act
	top := TopByMetricTypeAndYear(transactions, 3)

	// assert
	assert.Equal(t, expected, top)
}

// test fails
func TestNoData(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		// act
		top := TopByMetricTypeAndYear(nil, 10)
		// assert
		assert.Empty(t, top)
	})
	t.Run("empty slice", func(t *testing.T) {

		// act
		top := TopByMetricTypeAndYear([]*socrata.RealEstateTransaction{}, 10)
		// assert
		assert.Empty(t, top)
	})
}
