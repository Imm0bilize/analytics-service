package database

import "errors"

var (
	ErrParseConfigFile = errors.New("error when processing the config")
	ErrConnectionToDb  = errors.New("failed to connect to the database")
	ErrCreateQuery     = errors.New("failed to create query")
)
