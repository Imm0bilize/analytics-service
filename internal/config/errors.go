package config

import "errors"

var (
	ErrPathNotSpecified = errors.New("the path to the config is not specified")
	ErrGetFullPath      = errors.New("could not get the full path to the configuration file")
	ErrParseFile        = errors.New("error during file parsing")
	ErrOpenCfgFile      = errors.New("could not open the file")
	ErrReadCfgFile      = errors.New("file reading error")
	ErrCloseCfgFile     = errors.New("failed to close config file")
)
