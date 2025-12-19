package services

import "errors"

var ErrInvalidParams = errors.New("parameter validation failed")

var ErrInvalidUser = errors.New("invalid credentials")
var ErrUserLocked = errors.New("user locked")
var ErrUserDetailsUpdate = errors.New("user not affected")

var ErrInvalidProduct = errors.New("invalid credentials")
var ErrProductFetch = errors.New("product fetch")
var ErrProductDelete = errors.New("nothing changed")
