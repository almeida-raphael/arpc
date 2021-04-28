package errors


// Remote remote error struct on go error
type Remote struct {
	Err error
}

// Error compliance with error interface
func (e *Remote) Error() string {
	return e.Err.Error()
}