package repositories

import "errors"

var ErrRowsNotAffected = errors.New("no rows affected")
var ErrDBQuery = errors.New("db query error")
