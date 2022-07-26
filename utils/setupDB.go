package utils

import "github.com/NamanBalaji/keynetic/kv"

var Store *kv.DB

func InitStore() {
	Store = kv.NewDB()
}
