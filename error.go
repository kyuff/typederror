package typederror

import "fmt"

type enumConstraint interface {
	~string | ~int | ~int64 | ~int32
}

func Wrap[T enumConstraint](errorType T, err error) Error[T] {
	return Error[T]{
		err:  err,
		Type: errorType,
	}
}

type Error[T enumConstraint] struct {
	err  error
	Type T
}

func (err Error[T]) Error() string {
	return fmt.Sprintf("%v: %v", err.Type, err.err)
}

func (err Error[T]) Is(other error) bool {
	if other == nil {
		return false
	}
	a, ok := other.(Error[T])
	if !ok {
		return false
	}

	return a.Type == err.Type
}

func (err Error[T]) Unwrap() error {
	return err.err
}
