package apperrs

import "errors"

type AppError error

var ErrEntryNotFound AppError = errors.New("entry not found")
