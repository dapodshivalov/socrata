package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"socrata/src/socrata"
)

type FloatingTimestamp time.Time

// UnmarshalJSON for CustomTime to parse time in RFC3339 format
func (ct *FloatingTimestamp) UnmarshalJSON(b []byte) error {
	// Trim quotes, since JSON strings are quoted
	s := string(b)
	if s == "null" {
		return nil
	}
	s = strings.Trim(s, "\\\"")
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*ct = FloatingTimestamp(t)
	return nil
}

func (ct FloatingTimestamp) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02T15:04:05"))), nil
}

type FloatNumber float32

func (fl *FloatNumber) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	s = strings.Trim(s, "\\\"")
	parsedValue, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*fl = FloatNumber(parsedValue)
	return nil
}

func (fl FloatNumber) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%f\"", fl)), nil
}

type RawData struct {
	SerialNumber string            `json:"serialnumber"`
	Town         string            `json:"town"`
	SaleTs       FloatingTimestamp `json:"daterecorded"`
	SalesRatio   FloatNumber       `json:"salesratio"`
	SalesAmount  FloatNumber       `json:"saleamount"`
}

func (data *RawData) ToDomain() *socrata.RealEstateTransaction {
	if data == nil {
		return nil
	}
	return &socrata.RealEstateTransaction{
		Town:        socrata.Town(data.Town),
		SaleYear:    socrata.Year(time.Time(data.SaleTs).Year()),
		SalesRatio:  float64(data.SalesRatio),
		SalesVolume: float64(data.SalesAmount),
	}
}
