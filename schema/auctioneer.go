package schema

import (
	"errors"
	"strings"
)

// Auctioneer holds the basic details about the Auctioneer
type Auctioneer struct {
}

// AuctionReq ad request
type AuctionReq struct {
	AuctionID string `json:"auction_id"`
}

// Ok validates the request
func (b *AuctionReq) Ok() error {
	switch {
	case strings.TrimSpace(b.AuctionID) == "":
		return errors.New("auction id is required")
	}

	return nil
}
