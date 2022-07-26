package utils

import "os"

type environmnetVars struct {
	FwdAddr string
}

var Env environmnetVars

func InitEnvVars() {
	addr := os.Getenv("FORWARDING_ADDRESS")
	Env.FwdAddr = addr
}
