package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	IO_BLOCK_SIZE = "IO_BLOCK_SIZE"
	IO_MAX_BLOCKS = "IO_MAX_BLOCKS"
	IO_PATH       = "IO_PATH"
)

type Config struct {
	BlockSize int
	MaxBlocks int64
	Path      string
}

// Load is loading the configuration items from .env file.
func Load() (*Config, error) {

	c := Config{}
	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, "loading .env file")
	}

	val, defined := os.LookupEnv(IO_BLOCK_SIZE)
	if !defined {
		return nil, errors.New(fmt.Sprint("Could not read", IO_BLOCK_SIZE, "from config file"))
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Unable to use the", IO_BLOCK_SIZE, "config item value. Reason:", err))
	}
	c.BlockSize = n

	val, defined = os.LookupEnv(IO_MAX_BLOCKS)
	if !defined {
		return nil, errors.New(fmt.Sprint("Could not read", IO_MAX_BLOCKS, "from config file"))
	}
	n, err = strconv.Atoi(val)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Unable to use the", IO_MAX_BLOCKS, "config item value. Reason:", err))
	}
	c.MaxBlocks = int64(n)

	val, defined = os.LookupEnv(IO_PATH)
	if !defined {
		return nil, errors.New(fmt.Sprint("Could not read", IO_PATH, "from config file"))
	}
	c.Path = val

	return &c, nil
}
