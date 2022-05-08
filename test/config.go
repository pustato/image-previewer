package test

import (
	"net"
	"os"

	"code.cloudfoundry.org/bytefmt"
)

type Config struct {
	cacheSize   uint64
	serviceAddr string
	staticAddr  string
	cacheDir    string
}

func NewConfig() *Config {
	cacheSize, err := bytefmt.ToBytes(os.Getenv("CACHE_SIZE"))
	if err != nil {
		os.Exit(1)
	}

	return &Config{
		cacheSize:   cacheSize,
		serviceAddr: net.JoinHostPort(os.Getenv("HTTP_HOST"), os.Getenv("HTTP_PORT")),
		staticAddr:  os.Getenv("STATIC_HOST"),
		cacheDir:    os.Getenv("CACHE_DIR"),
	}
}
