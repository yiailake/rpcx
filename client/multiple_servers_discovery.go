package client

import (
	"time"

	"github.com/smallnest/rpcx/log"
)

// MultipleServersDiscovery is a multiple servers service discovery.
// It always returns the current servers and uses can change servers dynamically.
type MultipleServersDiscovery struct {
	pairs []*KVPair
	chans []chan []*KVPair
}

// NewMultipleServersDiscovery returns a new MultipleServersDiscovery.
func NewMultipleServersDiscovery(pairs []*KVPair) ServiceDiscovery {
	return &MultipleServersDiscovery{
		pairs: pairs,
	}
}

// GetServices returns the configured server
func (d MultipleServersDiscovery) GetServices() []*KVPair {
	return d.pairs
}

// WatchService returns a nil chan.
func (d *MultipleServersDiscovery) WatchService() chan []*KVPair {
	ch := make(chan []*KVPair, 10)
	d.chans = append(d.chans, ch)
	return ch
}

// Update is used to update servers at runtime.
func (d *MultipleServersDiscovery) Update(pairs []*KVPair) {
	for _, ch := range d.chans {
		ch := ch
		go func() {
			select {
			case ch <- pairs:
			case <-time.After(time.Minute):
				log.Warn("chan is full and new change has ben dropped")
			}
		}()
	}
}
