package localjson

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"socrata/src/socrata"
	"socrata/src/socrata/client/internal"
)

type Client struct {
	filepath string
}

func NewClient(filepath string) *Client {
	return &Client{filepath: filepath}
}

func (client *Client) GetTransactions() ([]*socrata.RealEstateTransaction, error) {
	file, err := os.Open(client.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var transactions []*internal.RawData
	if err := json.NewDecoder(file).Decode(&transactions); err != nil {
		return nil, err
	}
	result := make([]*socrata.RealEstateTransaction, 0, len(transactions))
	for _, tx := range transactions {
		if time.Time(tx.SaleTs).IsZero() {
			log.Printf("Empty year for tx: %s\n", tx.SerialNumber)
			continue
		}
		result = append(result, tx.ToDomain())
	}
	return result, nil
}
