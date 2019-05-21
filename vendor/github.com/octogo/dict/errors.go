package dict

import "errors"

var (
	errNotFound    = errors.New("key not found")
	errUnknownType = errors.New("can not consume unknown type")
)
