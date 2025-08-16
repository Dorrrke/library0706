package internal

import (
	"cmp"
	"flag"
	"os"
)

type Config struct {
	Addr        string
	Port        int
	DBDSN       string
	MigratePath string
	Debug       bool
}

const (
	defaultAddr  = "0.0.0.0"
	defaultPort  = 8080
	defaultDBStr = "postgresql://user:password@localhost:5432/library?sslmode=disable"
)

func ReadConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Addr, "addr", defaultAddr, "flag for use custom server addr")
	flag.IntVar(&cfg.Port, "port", defaultPort, "flag for use custom server port")
	flag.StringVar(&cfg.DBDSN, "db", defaultDBStr, "flag for setup db connection string")
	flag.StringVar(&cfg.MigratePath, "migrate", "migrations", "flag for setup path to migrations")
	flag.BoolVar(&cfg.Debug, "debug", false, "flag for enable debug mode")
	flag.Parse()

	cfg.DBDSN = cmp.Or(os.Getenv("DB_DSN"), cfg.DBDSN)
	cfg.MigratePath = cmp.Or(os.Getenv("MIGRATE_PATH"), cfg.MigratePath)

	return &cfg
}
