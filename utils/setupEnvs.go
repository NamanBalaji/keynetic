package utils

import (
	"errors"
	"os"
	"strings"
)

var Views []string
var SocketAddr string

func InitEnvVars() {
	views := os.Getenv("VIEWS")
	Views = strings.Split(views, ",")
	SocketAddr = os.Getenv("SOCKET_ADDRESS")
}

func RemoveFromView(v string) error {
	idx, ok := Contains(v)
	if !ok {
		return errors.New("view not present")
	}
	Views = append(Views[:idx], Views[idx+1:]...)
	return nil
}

func Contains(v string) (int, bool) {
	idx := -1
	for i, view := range Views {
		if view == v {
			idx = i
		}
	}
	if idx == -1 {
		return idx, false
	}

	return idx, true
}
