package utils

import (
	"github.com/NamanBalaji/keynetic/addresses"
	"github.com/NamanBalaji/keynetic/kv"
)

var Store *kv.DB
var View *addresses.Addrs

func InitStore() {
	Store = kv.NewDB()
}

func InitViews(v []string, sAddr string) {
	View = addresses.SetupAddrs(v, sAddr)
}
