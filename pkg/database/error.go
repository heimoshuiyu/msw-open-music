package database

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("object not found")
	ErrTagNotFound = errors.New("tag not found")
)
