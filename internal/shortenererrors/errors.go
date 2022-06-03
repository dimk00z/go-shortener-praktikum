package shortenererrors

import "errors"

var ErrURLDeleted = errors.New("shortURL has been deleted")

var ErrURLNotFound = errors.New("shortURL not found")
