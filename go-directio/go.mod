module github.com/devisions/go-playground/go-directio

go 1.15

require (
	github.com/joho/godotenv v1.3.0
	github.com/ncw/directio v1.0.5
	github.com/pkg/errors v0.9.1
)

// For now, dropped the usage of directio flavor that allows setting the `AlignSize`.
// replace github.com/ncw/directio v1.0.5 => /home/devisions/dev/devisions-gh/directio
