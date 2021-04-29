package errors

import "errors"

// Remote remote error struct on go error
type Remote struct {
	Err error
}

// Error compliance with error interface
func (e *Remote) Error() string {
	return e.Err.Error()
}

// Is method to comply with new errors functions
func (e *Remote) Is(target error) bool {
	tar, ok := target.(*Remote)
	if !ok {
		return false
	}

	return errors.Is(e.Err, tar.Err) || tar.Err == nil
}