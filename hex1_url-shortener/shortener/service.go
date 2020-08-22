package shortener

type ShortUrlService interface {
	Find(code string) (*ShortUrl, error)
	Store(redirect *ShortUrl) error
}
