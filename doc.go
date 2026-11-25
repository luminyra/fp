// Package fp provides three functional container types for Go 1.27:
// [Option] for a value that may be absent, [Result] for a value that may have
// failed, and [Either] for a value of one of two types.
//
// # Chaining across types
//
// Methods such as [Option.Map], [Option.AndThen], [Result.AndThen] and
// [Either.Fold] are generic methods, so a chain can change the value type
// without leaving the receiver:
//
//	port := fp.ResultOf(strconv.Atoi(raw)).
//		AndThen(checkPort).
//		Map(func(n int) string { return ":" + strconv.Itoa(n) })
//
// [Result.AndThen] takes an ordinary func(T) (R, error), so existing
// functions plug in unchanged, and [ResultOf] adapts any (value, error)
// return.
//
// # Interoperability
//
// [Option.All] exposes an Option as an [iter.Seq], [Option.Pointer] and
// [OptionFromPointer] convert to and from pointers, and Option implements
// [encoding/json.Marshaler] and [encoding/json.Unmarshaler], mapping an
// empty Option to null.
//
// # Guidelines
//
// All types are immutable values: methods return new values and never modify
// the receiver, so they are safe to share between goroutines as long as the
// wrapped values are. Prefer plain (T, bool) and (T, error) returns in public
// APIs, and use these types where chaining makes the code clearer.
//
// The package is inspired by samber/mo and by Rust's Option and Result.
package fp
