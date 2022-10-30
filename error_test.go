package typederror_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/kyuff/typederror"
)

type DomainError string

const (
	InvalidArgument      DomainError = "INVALID_ARGUMENT"
	ConstraintViolation  DomainError = "CONSTRAINT_VIOLATION"
	InsufficientCapacity DomainError = "INSUFFICIENT_CAPACITY"
)

type SystemError int

const (
	ErrorA SystemError = iota
	ErrorB
	ErrorC
)

func (err SystemError) String() string {
	switch err {
	case ErrorA:
		return "a"
	case ErrorB:
		return "b"
	case ErrorC:
		return "c"
	default:
		return "unknown"
	}
}

type testConstraint interface {
	~string | ~int
}

func assertTypedErrorWithType[T testConstraint](t *testing.T, got error, expected T) {
	var typedError typederror.Error[T]
	if !errors.As(got, &typedError) {
		t.FailNow()
	}
	if typedError.Type != expected {
		t.FailNow()
	}
}

func TestWrap(t *testing.T) {

	t.Run("should support errors.As directly for string enums", func(t *testing.T) {
		// act
		err := typederror.Wrap(InvalidArgument, io.EOF)

		// assert
		assertTypedErrorWithType(t, err, InvalidArgument)
	})

	t.Run("should support errors.As indirectly for string enums", func(t *testing.T) {
		// act
		err := fmt.Errorf("test error: %w", typederror.Wrap(InsufficientCapacity, io.EOF))

		// assert
		assertTypedErrorWithType(t, err, InsufficientCapacity)
	})

	t.Run("should support errors.As directly for int enums", func(t *testing.T) {
		// act
		err := typederror.Wrap(ErrorB, io.EOF)

		// assert
		assertTypedErrorWithType(t, err, ErrorB)
	})

	t.Run("should support errors.As indirectly for int enums", func(t *testing.T) {
		// act
		err := fmt.Errorf("test error: %w", typederror.Wrap(ErrorA, io.EOF))

		// assert
		assertTypedErrorWithType(t, err, ErrorA)
	})

	t.Run("should give the correct error string for string enums", func(t *testing.T) {
		// arrange
		err := typederror.Wrap(ConstraintViolation, io.EOF)

		// act
		got := err.Error()

		if got != "CONSTRAINT_VIOLATION: EOF" {
			t.FailNow()
		}
	})

	t.Run("should give the correct error string for int enums", func(t *testing.T) {
		// arrange
		err := typederror.Wrap(ErrorB, io.EOF)

		// act
		got := err.Error()

		if got != "b: EOF" {
			t.FailNow()
		}
	})

	t.Run("should support Unwrap directly for string enums", func(t *testing.T) {
		// act
		err := typederror.Wrap(InvalidArgument, io.EOF)

		// assert
		if !errors.Is(err, io.EOF) {
			t.FailNow()
		}
	})

	t.Run("should support Unwrap directly for int enums", func(t *testing.T) {
		// act
		err := typederror.Wrap(InvalidArgument, io.EOF)

		// assert
		if !errors.Is(err, io.EOF) {
			t.FailNow()
		}
	})

	t.Run("should support errors.Is directly for string enums", func(t *testing.T) {
		var (
			err      = typederror.Wrap(ConstraintViolation, io.EOF)
			expected = typederror.Wrap(ConstraintViolation, io.EOF)
			other    = typederror.Wrap(InvalidArgument, io.EOF)
		)

		if !errors.Is(err, expected) {
			t.FailNow()
		}
		if errors.Is(err, other) {
			t.FailNow()
		}
	})

	t.Run("should support errors.Is directly for int enums", func(t *testing.T) {
		var (
			err      = typederror.Wrap(ErrorC, io.EOF)
			expected = typederror.Wrap(ErrorC, io.EOF)
			other    = typederror.Wrap(ErrorB, io.EOF)
		)

		if !errors.Is(err, expected) {
			t.FailNow()
		}
		if errors.Is(err, other) {
			t.FailNow()
		}
	})
}
