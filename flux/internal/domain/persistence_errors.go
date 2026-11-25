package domain

// PersistenceTemporaryError represents a persistence failure
// that may succeed on a subsequent attempt.
type PersistenceTemporaryError struct {
	Err error
}

func (e *PersistenceTemporaryError) Error() string {
	return "temporary persistence error: " + e.Err.Error()
}

func (e *PersistenceTemporaryError) Unwrap() error {
	return e.Err
}

// PersistencePermanentError represents a persistence failure
// that requires correction before retrying.
type PersistencePermanentError struct {
	Err error
}

func (e *PersistencePermanentError) Error() string {
	return "permanent persistence error: " + e.Err.Error()
}

func (e *PersistencePermanentError) Unwrap() error {
	return e.Err
}
