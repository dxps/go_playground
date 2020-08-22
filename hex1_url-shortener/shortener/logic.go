package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrShortUrlNotFound = errors.New("ShortUrl Not Found")
	ErrShortUrlInvalid  = errors.New("ShortUrl Invalid")
)

type shortUrlService struct {
	repo ShortUrlRepository
}

func NewShortUrlService(repo ShortUrlRepository) ShortUrlService {
	return &shortUrlService{
		repo,
	}
}

func (s *shortUrlService) Find(code string) (*ShortUrl, error) {
	return s.repo.Find(code)
}

func (s *shortUrlService) Store(su *ShortUrl) error {
	if err := validate.Validate(su); err != nil {
		return errs.Wrap(ErrShortUrlInvalid, "service.ShortUrl.Store")
	}
	su.Code = shortid.MustGenerate()
	su.CreatedAt = time.Now().UTC().Unix()
	return s.repo.Store(su)
}
