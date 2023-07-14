package utils

import (
	"github.com/NamanBalaji/keynetic/addresses"
	"github.com/NamanBalaji/keynetic/kv"
	"github.com/NamanBalaji/keynetic/shard"
	vectorclock "github.com/NamanBalaji/keynetic/vectorClock"
)

var Store *kv.DB
var View *addresses.Addrs
var Vc vectorclock.VectorClock
var Shard *shard.Shard

func InitStore() {
	Store = kv.NewDB()
}

func SetStore(newStore map[string]string) {
	if newStore != nil {
		Store.Database = newStore
	}

}

func InitViews(v []string, sAddr string) {
	View = addresses.SetupAddrs(v, sAddr)
}

func SetVectorClock(vc map[string]int) {
	for key, val := range Vc {
		Vc[key] = max(val, vc[key])
	}
}

func InitVectorClock(views []string) {
	Vc = vectorclock.NewVectorClock(views)
}

func InitShard(shardCount int, socketAddr string) {
	Shard = shard.NewShard(shardCount, socketAddr, View.Views)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
