package addresses

import (
	"errors"
	"fmt"
)

type Addrs struct {
	Views      []string
	SocketAddr string
}

func SetupAddrs(v []string, addr string) *Addrs {
	return &Addrs{
		Views:      v,
		SocketAddr: addr,
	}
}

// Contains checks for the presence of a given address string in views slice
func (a *Addrs) Contains(v string) (int, bool) {
	for i, view := range a.Views {
		if v == view {
			return i, true
		}
	}
	return -1, false
}

// RemoveFromView removes a given address from the view
func (a *Addrs) RemoveFromView(v string) error {
	if v == a.SocketAddr {
		return errors.New("cannot remove self from view")
	}
	idx := -1
	for i, view := range a.Views {
		if v == view {
			idx = i
		}
	}
	if idx == -1 {
		return fmt.Errorf("address: %s, not present in view", v)
	}
	a.Views = append(a.Views[:idx], a.Views[idx+1:]...)
	return nil
}

// AddToView adds a given address to view
func (a *Addrs) AddToView(v string) {
	a.Views = append(a.Views, v)
}
