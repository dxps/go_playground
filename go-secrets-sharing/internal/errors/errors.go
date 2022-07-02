package errors

import "errors"

type AppError error

var EntryNotFound AppError = errors.New("entry not found")
