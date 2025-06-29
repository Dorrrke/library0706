package internal

import "flag"

type Config struct {
	Addr string
	Port int
}

const (
	defaultAddr = "0.0.0.0"
	defaultPort = 8080
)

func ReadConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Addr, "addr", defaultAddr, "flag for use custom server addr")
	flag.IntVar(&cfg.Port, "port", defaultPort, "flag for use custom server port")
	flag.Parse()

	return &cfg
}
