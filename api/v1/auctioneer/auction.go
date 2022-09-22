package auctioneer

import (
	pkg "bidding/pkg"
	"bidding/schema"
	"bidding/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type BidResponse struct {
	BidderID string  `json:"bidder_id"`
	Amount   float64 `json:"amount"`
}

type bidds []*BidResponse

func (a bidds) Len() int           { return len(a) }
func (a bidds) Less(i, j int) bool { return a[i].Amount > a[j].Amount }
func (a bidds) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// - gets all the bidders
// - start a bid section by making req to all the bidders with ad auction details
// - collects the response from the bidders, anounce the winner
func auctionHandler(w http.ResponseWriter, r *http.Request) *string {
	var (
		input schema.AuctionReq
		wg    sync.WaitGroup
	)
	s := ""
	if err := utils.Decode(r, &input); err != nil {
		s = "got error whle decoding"
		return &s
	}
	bidders := store.Bidder().List()
	if len(bidders) == 0 {
		s = "no bidders available for this auction"
		return &s
	}
	wg.Add(len(bidders))

	data := make(chan *BidResponse, len(bidders))
	for _, b := range bidders {
		go collectBidResponse(input.AuctionID, b.Host, &wg, data)
	}

	var bidRes bidds
	for i := 0; i < len(bidders); i++ {
		if d := <-data; d != nil {
			bidRes = append(bidRes, d)
		}
	}

	wg.Wait()
	close(data)
	sort.Sort(bidRes)
	if len(bidRes) == 0 {
		s = "bidders not responding with in time"
		return &s
	}

	pkg.OK(w, bidRes[0])
	return nil
}

func collectBidResponse(auctionID, host string, wg *sync.WaitGroup,
	data chan *BidResponse) {
	var err error
	body := bytes.NewBuffer(nil)
	json.NewEncoder(body).Encode(map[string]interface{}{
		"auction_id": auctionID,
	})

	defer func() {
		wg.Done()
		if err != nil {
			fmt.Println("null data, ", err.Error())
			data <- nil
		}
	}()

	url := host + "/v1/bid"
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 190*time.Millisecond)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var res struct {
		Data *BidResponse `json:"data"`
		Meta pkg.Meta     `json:"meta"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}

	data <- res.Data
}
