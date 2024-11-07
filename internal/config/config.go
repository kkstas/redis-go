package config

import (
	"flag"
)

type Config struct {
	Dir        string
	DBFilename string
}

func New() *Config {
	dir := flag.String("dir", "", "The directory where Redis files are stored")
	dbFilename := flag.String("dbfilename", "", "The name of the Redis dump file")
	flag.Parse()

	return &Config{
		Dir:        *dir,
		DBFilename: *dbFilename,
	}
}
