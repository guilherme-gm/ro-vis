package domain

type NotFoundError struct {
	message string
}

func (e NotFoundError) Error() string {
	return e.message
}

func (e NotFoundError) Is(target error) bool {
	_, ok := target.(NotFoundError)
	return ok
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{message: message}
}
