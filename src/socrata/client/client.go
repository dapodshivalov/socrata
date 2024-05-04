package client

import (
	"encoding/json"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"socrata/src/socrata"
	"socrata/src/socrata/client/internal"

	"golang.org/x/sync/errgroup"

	"github.com/SebastiaanKlippert/go-soda"
)

const (
	batchSize = 5000
	logRate   = 20 // log every `logRate` loaded batch in order to see progress in console
)

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
	count := uint32(0)

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
				atomic.AddUint32(&count, uint32(len(data)))
				if count%(batchSize*logRate) == 0 {
					log.Printf("got %d/%d records", count, offsetRequest.Count())
				}

				dataChannel <- data
			}
			return nil
		})
	}
	go func() {
		g.Wait()
		close(dataChannel)
	}()
	for data := range dataChannel {
		for _, d := range data {
			if time.Time(d.SaleTs).IsZero() {
				log.Printf("Empty year for tx: %s\n", d.SerialNumber)
				continue
			}
			result = append(result, d.ToDomain())
		}
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func (client *Client) makeOffsetRequest() (*soda.OffsetGetRequest, error) {
	req := soda.NewGetRequest("https://data.ct.gov/resource/5mzw-sjtu", client.token)
	req.Query.Select = []string{"daterecorded", "town", "salesratio", "saleamount", "serialnumber"}
	req.Format = "json"
	req.Query.AddOrder("daterecorded", soda.DirAsc)
	return soda.NewOffsetGetRequest(req)
}
