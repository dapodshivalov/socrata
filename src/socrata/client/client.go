package client

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"time"

	"socrata/src/socrata"
	"socrata/src/socrata/client/internal"

	"golang.org/x/sync/errgroup"

	"github.com/SebastiaanKlippert/go-soda"
)

const batchSize = 5000

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (client *Client) GetTransactions() ([]*socrata.RealEstateTransaction, error) {
	result := make([]*socrata.RealEstateTransaction, 0, batchSize)

	offsetRequest, err := client.makeOffsetRequest()
	if err != nil {
		return nil, err
	}

	dataChannel := make(chan []*internal.RawData, 10)

	g := new(errgroup.Group)
	g.SetLimit(10)
	total := uint32(0)

	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for {
				resp, err := offsetRequest.Next(batchSize)
				if errors.Is(err, soda.ErrDone) {
					break
				}
				if err != nil {
					return err
				}

				var data []*internal.RawData
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					return err
				}
				atomic.AddUint32(&total, uint32(len(data)))
				log.Printf("goroutine %d got %d records. total %d", i, len(data), total)

				dataChannel <- data
			}
			return nil
		})
	}
	go func() {
		g.Wait()
		close(dataChannel)
	}()
	rawData := make([]*internal.RawData, 0, batchSize)
	for data := range dataChannel {
		rawData = append(rawData, data...)
		for _, d := range data {
			if time.Time(d.SaleTs).IsZero() {
				log.Printf("Empty year for tx: %s\n", d.SerialNumber)
				continue
			}
			result = append(result, d.ToDomain())
		}
		//result = append(result, util.Transform(data, (*internal.RawData).ToDomain)...)
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	//save raw data to file
	file, err := os.Create("dump.json")
	if err != nil {
		return nil, err
	}
	if err := json.NewEncoder(file).Encode(rawData); err != nil {
		return nil, err
	}

	log.Printf("Total transactions: %d\n", len(result))
	emptyYears := 0
	for _, transaction := range result {
		if transaction.SaleYear == 0 {
			emptyYears++
		}
	}
	log.Printf("Empty years: %d\n", emptyYears)
	return result, nil
}

func (client *Client) makeOffsetRequest() (*soda.OffsetGetRequest, error) {
	req := soda.NewGetRequest("https://data.ct.gov/resource/5mzw-sjtu", client.token)
	req.Query.Select = []string{"daterecorded", "town", "salesratio", "saleamount", "serialnumber"}
	req.Format = "json"
	req.Query.AddOrder("daterecorded", soda.DirAsc)
	return soda.NewOffsetGetRequest(req)
}
