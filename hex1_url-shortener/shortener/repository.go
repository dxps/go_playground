package shortener

type ShortUrlRepository interface {
	Find(code string) (*ShortUrl, error)
	Store(redirect *ShortUrl) error
}
