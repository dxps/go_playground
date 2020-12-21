package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	IO_BLOCK_SIZE    = "IO_BLOCK_SIZE"
	IO_PATH          = "IO_PATH"
	IO_FILE_MAX_SIZE = "IO_FILE_MAX_SIZE"
)

type Config struct {
	BlockSize   int
	Path        string
	FileMaxSize int64
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

	val, defined = os.LookupEnv(IO_PATH)
	if !defined {
		return nil, errors.New(fmt.Sprint("Could not read", IO_PATH, "from config file"))
	}
	c.Path = val

	val, defined = os.LookupEnv(IO_FILE_MAX_SIZE)
	if !defined {
		return nil, errors.New(fmt.Sprint("Could not read", IO_FILE_MAX_SIZE, "from config file"))
	}
	n, err = strconv.Atoi(val)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Unable to use the", IO_FILE_MAX_SIZE, "config item value. Reason:", err))
	}
	c.FileMaxSize = int64(n)

	return &c, nil
}
