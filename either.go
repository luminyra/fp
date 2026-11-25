package fp

import "fmt"

// Either holds a value of type L (left) or of type R (right). The zero Either
// is a Left holding the zero value of L.
type Either[L, R any] struct {
	left    L
	right   R
	isRight bool
}

// Left returns an Either holding value on the left side.
func Left[L, R any](value L) Either[L, R] {
	return Either[L, R]{left: value}
}

// Right returns an Either holding value on the right side.
func Right[L, R any](value R) Either[L, R] {
	return Either[L, R]{right: value, isRight: true}
}

// IsLeft reports whether e holds a left value.
func (e Either[L, R]) IsLeft() bool {
	return !e.isRight
}

// IsRight reports whether e holds a right value.
func (e Either[L, R]) IsRight() bool {
	return e.isRight
}

// Left returns the left value as an Option.
func (e Either[L, R]) Left() Option[L] {
	return OptionOf(e.left, !e.isRight)
}

// Right returns the right value as an Option.
func (e Either[L, R]) Right() Option[R] {
	return OptionOf(e.right, e.isRight)
}

// Swap returns an Either with the sides exchanged.
func (e Either[L, R]) Swap() Either[R, L] {
	if e.isRight {
		return Left[R, L](e.right)
	}
	return Right[R](e.left)
}

// Fold reduces e to a single value by calling onLeft or onRight, depending on
// the side it holds.
func (e Either[L, R]) Fold[U any](onLeft func(L) U, onRight func(R) U) U {
	if e.isRight {
		return onRight(e.right)
	}
	return onLeft(e.left)
}

// MapLeft transforms a left value with f; a right value is carried over.
func (e Either[L, R]) MapLeft[U any](f func(L) U) Either[U, R] {
	if e.isRight {
		return Right[U](e.right)
	}
	return Left[U, R](f(e.left))
}

// MapRight transforms a right value with f; a left value is carried over.
func (e Either[L, R]) MapRight[U any](f func(R) U) Either[L, U] {
	if e.isRight {
		return Right[L](f(e.right))
	}
	return Left[L, U](e.left)
}

// String formats e as "Left(value)" or "Right(value)".
func (e Either[L, R]) String() string {
	if e.isRight {
		return fmt.Sprintf("Right(%v)", e.right)
	}
	return fmt.Sprintf("Left(%v)", e.left)
}
