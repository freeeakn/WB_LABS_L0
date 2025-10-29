package order

import "fmt"

type MissingFieldError string

func (m MissingFieldError) Error() string { return fmt.Sprintf("missing field: %s", string(m)) }

func ErrMissingField(name string) error { return MissingFieldError(name) }

var ErrInvalidOrder = fmt.Errorf("invalid order")
