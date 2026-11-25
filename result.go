package fp

import (
	"errors"
	"fmt"
)

// errNilError replaces a nil error passed to Err so that the Result still
// reports failure.
var errNilError = errors.New("fp: Err called with nil error")

// Result holds either a value of type T or an error. The zero Result holds the
// zero value of T and no error.
type Result[T any] struct {
	value T
	err   error
}

// Ok returns a successful Result holding value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err returns a failed Result holding err. A nil err is replaced by a non-nil
// error so that the Result still reports failure.
func Err[T any](err error) Result[T] {
	if err == nil {
		err = errNilError
	}
	return Result[T]{err: err}
}

// Errf returns a failed Result whose error is fmt.Errorf(format, args...). The
// %w verb is supported.
func Errf[T any](format string, args ...any) Result[T] {
	return Err[T](fmt.Errorf(format, args...))
}

// ResultOf adapts a (value, error) pair, so that a call such as
// ResultOf(strconv.Atoi(s)) becomes a Result.
func ResultOf[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// IsOk reports whether r holds a value.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr reports whether r holds an error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Err returns the held error, or nil on success.
func (r Result[T]) Err() error {
	return r.err
}

// Get returns the value and the error, like a regular Go function.
func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

// MustGet returns the value and panics with the held error on failure.
func (r Result[T]) MustGet() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// Or returns the value, or fallback on failure.
func (r Result[T]) Or(fallback T) T {
	if r.err != nil {
		return fallback
	}
	return r.value
}

// OrZero returns the value, or the zero value of T on failure.
func (r Result[T]) OrZero() T {
	if r.err != nil {
		var zero T
		return zero
	}
	return r.value
}

// Option returns Some(value) on success and None on failure, dropping the
// error.
func (r Result[T]) Option() Option[T] {
	if r.err != nil {
		return None[T]()
	}
	return Some(r.value)
}

// Map returns Ok(f(value)) on success; a failure is carried over with the
// same error. The result type may differ from T.
func (r Result[T]) Map[R any](f func(T) R) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return Ok(f(r.value))
}

// AndThen calls f with the value on success and turns its (value, error)
// return into a Result; a failure is carried over and f is not called. f has
// the shape of an ordinary Go function, so existing functions plug in
// directly.
func (r Result[T]) AndThen[R any](f func(T) (R, error)) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return ResultOf(f(r.value))
}

// MapErr returns r with its error replaced by f(err), for example to add
// context with [fmt.Errorf]. A successful Result is returned unchanged. A nil
// return from f is treated like a nil error passed to Err.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.err == nil {
		return r
	}
	return Err[T](f(r.err))
}

// Recover calls f with the error on failure, so that it can produce a value
// or a different error; a successful Result is returned unchanged.
func (r Result[T]) Recover(f func(error) (T, error)) Result[T] {
	if r.err == nil {
		return r
	}
	return ResultOf(f(r.err))
}

// String formats r as "Ok(value)" or "Err(error)".
func (r Result[T]) String() string {
	if r.err != nil {
		return fmt.Sprintf("Err(%v)", r.err)
	}
	return fmt.Sprintf("Ok(%v)", r.value)
}
