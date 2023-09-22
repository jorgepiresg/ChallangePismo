package main

import (
	"github.com/jorgepiresg/ChallangePismo/config"
	"github.com/jorgepiresg/ChallangePismo/server"
)

func main() {
	cfg := config.New()
	server := server.New(cfg)
	server.Start()
}
