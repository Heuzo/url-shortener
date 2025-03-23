package storage

import "errors"

var (
	ErrURLNotFound   = errors.New("url not found")
	ErrURLExists     = errors.New("url already exists")
	ErrAliasNotFound = errors.New("alias not found")
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SQLService
type SQLService interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
	AliasExists(alias string) (bool, error)
}
