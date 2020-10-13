package shortener

type ShortUrlSerializer interface {
	Decode(input []byte) (*ShortUrl, error)
	Encode(input *ShortUrl) ([]byte, error)
}
