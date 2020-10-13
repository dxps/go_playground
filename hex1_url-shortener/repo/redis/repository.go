package redis

import (
	"fmt"
	"strconv"

	"devisions.org/go-playground/hex1_url-shortener/shortener"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepo struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.ShortUrlRepository, error) {
	repo := &redisRepo{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepo) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepo) Find(code string) (*shortener.ShortUrl, error) {
	su := &shortener.ShortUrl{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.ShortUrl.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrShortUrlNotFound, "repository.ShortUrl.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.ShortUrl.Find")
	}
	su.Code = data["code"]
	su.URL = data["url"]
	su.CreatedAt = createdAt
	return su, nil
}

func (r *redisRepo) Store(su *shortener.ShortUrl) error {
	key := r.generateKey(su.Code)
	data := map[string]interface{}{
		"code":       su.Code,
		"url":        su.URL,
		"created_at": su.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.ShortUrl.Store")
	}
	return nil
}
