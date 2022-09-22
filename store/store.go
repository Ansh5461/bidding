package store

import (
	"bidding/schema"
)

var bidders = make([]*schema.Bidder, 0)

type Conn struct {
	BidderConn Bidders
}

type BidderStore struct {
	*Conn
}

type Store interface {
	Bidder() Bidders
}

type Bidders interface {
	Add(bidder *schema.Bidder)
	List() []*schema.Bidder
	Count() int
}

func NewBidderStore(st *Conn) *BidderStore {
	return &BidderStore{st}
}

func (b *BidderStore) Add(bidder *schema.Bidder) {
	bidders = append(bidders, bidder)
}

func (b *BidderStore) List() []*schema.Bidder {
	return bidders
}

func (b *BidderStore) Count() int {
	return len(bidders)
}

func NewStore() *Conn {
	conn := new(Conn)
	conn.BidderConn = NewBidderStore(conn)

	return conn
}

func (s *Conn) Bidder() Bidders {
	return s.BidderConn
}
