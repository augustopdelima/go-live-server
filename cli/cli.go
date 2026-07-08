package cli

import (
	"flag"
)

type Config struct {
	Port      int
	Directory string
}

func Parse() Config {
	config := Config{}

	flag.StringVar(&config.Directory, "dir", "./", "directory")

	flag.IntVar(&config.Port, "port", 5000, "port http server")

	flag.Parse()

	return config
}
