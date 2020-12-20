package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	IO_BLOCKSIZE = "IO_BLOCKSIZE"
	IO_FILEPATH  = "IO_FILEPATH"
)

type Config struct {
	BlockSize int
	Filepath  string
}

// Loading config from .env file.
func Load() (*Config, error) {
	c := Config{}
	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, "loading .env file")
	}
	if val, defined := os.LookupEnv(IO_BLOCKSIZE); defined {
		if n, err := strconv.Atoi(val); err != nil {
			log.Fatal("Unable to use the", IO_BLOCKSIZE, "config item value. Reason:", err)
		} else {
			c.BlockSize = n
		}
	}
	if val, defined := os.LookupEnv(IO_FILEPATH); defined {
		c.Filepath = val
	}
	return &c, nil
}
