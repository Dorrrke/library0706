package errors

import "errors"

var (
	ErrUserAlredyExist = errors.New("user alredy exist")
	ErrIvalidCreds     = errors.New("ivalid email or pass")
)
