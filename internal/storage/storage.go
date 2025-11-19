package storage

import (
	"errors"
)

var ErrURLNotFound = errors.New("url not found")
var ErrURLExists = errors.New("url already exists")
