package json

import (
	"encoding/json"

	"devisions.org/go-playground/hex1_url-shortener/shortener"
	"github.com/pkg/errors"
)

type ShortUrl struct{}

func (s *ShortUrl) Decode(input []byte) (*shortener.ShortUrl, error) {
	su := &shortener.ShortUrl{}
	if err := json.Unmarshal(input, su); err != nil {
		return nil, errors.Wrap(err, "serializer.ShortUrl.Decode")
	}
	return su, nil
}

func (s *ShortUrl) Encode(input *shortener.ShortUrl) ([]byte, error) {
	suJson, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.ShortUrl.Encode")
	}
	return suJson, nil
}
