package fp

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
)

// ErrNone is the panic value of [Option.MustGet] on an empty Option.
var ErrNone = errors.New("fp: option is empty")

// Option holds either a value of type T or nothing. The zero Option is empty.
type Option[T any] struct {
	value  T
	isSome bool
}

// Some returns an Option holding value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, isSome: true}
}

// None returns an empty Option.
func None[T any]() Option[T] {
	return Option[T]{}
}

// OptionOf adapts the comma-ok idiom: it returns Some(value) when ok is true
// and None otherwise.
func OptionOf[T any](value T, ok bool) Option[T] {
	if !ok {
		return None[T]()
	}
	return Some(value)
}

// OptionFromPointer returns None for a nil pointer and Some of the pointed-to
// value otherwise.
func OptionFromPointer[T any](p *T) Option[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

// IsSome reports whether o holds a value.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// IsNone reports whether o is empty.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// Get returns the value and whether o holds one.
func (o Option[T]) Get() (T, bool) {
	return o.value, o.isSome
}

// MustGet returns the value and panics with ErrNone when o is empty.
func (o Option[T]) MustGet() T {
	if !o.isSome {
		panic(ErrNone)
	}
	return o.value
}

// Or returns the value, or fallback when o is empty.
func (o Option[T]) Or(fallback T) T {
	if !o.isSome {
		return fallback
	}
	return o.value
}

// OrFunc returns the value, or the result of fallback when o is empty.
// fallback is called only when needed.
func (o Option[T]) OrFunc(fallback func() T) T {
	if !o.isSome {
		return fallback()
	}
	return o.value
}

// OrZero returns the value, or the zero value of T when o is empty.
func (o Option[T]) OrZero() T {
	return o.value
}

// Pointer returns a pointer to a copy of the value, or nil when o is empty.
func (o Option[T]) Pointer() *T {
	if !o.isSome {
		return nil
	}
	return new(o.value)
}

// Map returns Some(f(value)) when o holds a value and None otherwise. The
// result type may differ from T.
func (o Option[T]) Map[R any](f func(T) R) Option[R] {
	if !o.isSome {
		return None[R]()
	}
	return Some(f(o.value))
}

// AndThen returns f(value) when o holds a value and None otherwise. Use it to
// chain steps that may themselves produce nothing.
func (o Option[T]) AndThen[R any](f func(T) Option[R]) Option[R] {
	if !o.isSome {
		return None[R]()
	}
	return f(o.value)
}

// Filter keeps the value only when keep reports true.
func (o Option[T]) Filter(keep func(T) bool) Option[T] {
	if !o.isSome || !keep(o.value) {
		return None[T]()
	}
	return o
}

// All returns an iterator that yields the value once when o holds one, so
// that an Option can be ranged over or passed to iterator-based APIs.
func (o Option[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.isSome {
			yield(o.value)
		}
	}
}

// String formats o as "Some(value)" or "None".
func (o Option[T]) String() string {
	if !o.isSome {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", o.value)
}

// MarshalJSON encodes the value as its JSON form and an empty Option as null.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.isSome {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON decodes null as None and any other JSON value as Some.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = None[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("fp: decode option: %w", err)
	}
	*o = Some(v)
	return nil
}
