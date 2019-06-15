package gohellomod

import (
	quoteV3 "rsc.io/quote/v3"
)

// Hello returns a random quote.
func Hello() string {
	return quoteV3.HelloV3()
}

// Proverb returns a concurrency related quote.
func Proverb() string {
	return quoteV3.Concurrency()
}
