package postgre

import "errors"

var (
	ErrParseConfigFile = errors.New("error when processing the config")
	ErrConnectionToDb  = errors.New("failed to connect to the database")
)
