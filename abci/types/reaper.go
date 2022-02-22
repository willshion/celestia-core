package types

// Reaper provides applications more fine grained access to the mempool by
// retrieving individual transactions at a time.
type Reaper interface {
	Next() (tx []byte, empty bool)
}

// NewInProcessReaper creates a Reaper that is designed to work for in process
// ABCI applications
func NewInProcessReaper() Reaper {
	return &inProcessReaper{}
}

var _ Reaper = &inProcessReaper{}

type inProcessReaper struct {
}

func (ipr *inProcessReaper) Next() (tx []byte, empty bool) {
	return nil, true
}
