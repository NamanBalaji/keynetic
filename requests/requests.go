package requests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NamanBalaji/keynetic/utils"
)

// SendHeartbeats sends ping request to each replica and returns a list of down instances
func SendHeartbeats() []string {
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()

	var down []string
	for _, v := range utils.View.Views {
		if v != utils.View.SocketAddr {
			url := fmt.Sprintf("http://%s/ping", v)
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			_, err := http.DefaultClient.Do(req)
			if err != nil {
				down = append(down, v)
			}
		}
	}
	return down
}

// BroadcastDeleteView sends a request to delete an addresses from a replica's view
func BroadcastDeleteView(addrSend string, addrDelete ...string) (string, error) {
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()

	for _, v := range addrDelete {
		url := fmt.Sprintf("http://%s/broadcast-delete/%s", addrSend, v)
		req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("replica %s is down or busy", addrSend)
			return addrSend, err
		}
	}
	return "", nil
}

// BroadcastPutView sends a request to add a given address to a replica's view
func BroadcastPutView(addrSend, addrPut string) error {
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()
	url := fmt.Sprintf("http://%s/broadcast-put/%s", addrSend, addrPut)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	_, err := http.DefaultClient.Do(req)
	return err
}

// GetKeyValueStore sends a get request to fetch the local key-value store of a given replica
func GetKeyValueStore(addr string) (*http.Response, error) {
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()

	url := fmt.Sprintf("http://%s/store", addr)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	return http.DefaultClient.Do(req)
}
