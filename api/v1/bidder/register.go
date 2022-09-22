package bidder

import (
	pkg "bidding/pkg"
	"bidding/schema"
	"bidding/utils"
	"fmt"
	"net/http"
)

func getBiddersHandler(w http.ResponseWriter, r *http.Request) *string {
	pkg.OK(w, store.Bidder().List())
	s := ""
	return &s
}

func bidderRegisterHandler(w http.ResponseWriter, r *http.Request) *string {
	var req schema.BidderReq

	if err := utils.Decode(r, &req); err != nil {
		s := "got error in bidderRegisterHandler"
		return &s
	}

	// NOTE: considering the bidder request come from localhost
	bidder := &schema.Bidder{
		ID:    fmt.Sprintf("bidder_%d", store.Bidder().Count()+1),
		Name:  req.Name,
		Host:  fmt.Sprintf("http://%s", r.Host),
		Delay: req.Delay,
	}
	store.Bidder().Add(bidder)

	pkg.Created(w, bidder)
	return nil
}
