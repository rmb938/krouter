package storage

import "errors"

var (
	ErrExists   = errors.New("resource already exists")
	ErrNotFound = errors.New("resource is not found")
)
