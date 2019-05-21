package rate

import "errors"

var errTooFast = errors.New("too fast")

// IsTooFast returns true if the given error is the errTooFast.
func IsTooFast(err error) bool {
	return err == errTooFast
}
