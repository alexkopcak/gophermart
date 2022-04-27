package auth

import "errors"

var (
	ErrUserAlreadyExsist = errors.New("user already exsist")
	ErrInternalServer    = errors.New("internal server error")
	ErrBadLoginPassword  = errors.New("bad login/password")
	ErrUserNotExsist     = errors.New("user not exsist")
)
