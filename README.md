# fp

[![Go Reference](https://pkg.go.dev/badge/github.com/luminyra/fp.svg)](https://pkg.go.dev/github.com/luminyra/fp)
[![CI](https://github.com/luminyra/fp/actions/workflows/ci.yml/badge.svg)](https://github.com/luminyra/fp/actions/workflows/ci.yml)

Option, Result and Either for Go 1.27, chainable across types with generic methods.

```go
// checkPort is an ordinary func(int) (int, error).
addr := fp.ResultOf(strconv.Atoi(os.Getenv("PORT"))).
	AndThen(checkPort).
	Map(func(n int) string { return ":" + strconv.Itoa(n) }). // Result[int] becomes Result[string]
	Or(":8080")
```

## Install

```sh
go get github.com/luminyra/fp
```

Requires Go 1.27 or later. The module has no dependencies.

## Highlights

- **Generic methods.** `Map`, `AndThen`, `Fold`, `MapLeft` and `MapRight` change the value type right on the receiver, so chains never fall back to package-level helpers.
- **Plain Go functions plug in.** `Result.AndThen` takes an ordinary `func(T) (R, error)`, and `ResultOf` adapts any `(value, error)` return, so existing code needs no wrappers.
- **Works with the rest of Go.** `Option.All` is an `iter.Seq`, `Option.Pointer` and `OptionFromPointer` convert to and from pointers, and `Option` encodes to JSON as the value or `null`.
- **Immutable values.** Methods never modify the receiver, so the types are safe to share between goroutines.

## Usage

```go
// Option: a value that may be absent.
host := fp.OptionOf(os.LookupEnv("HOST")).Or("localhost")

// Result: a value that may have failed.
cfg := fp.ResultOf(os.ReadFile(path)).
	AndThen(parseConfig).
	MapErr(func(err error) error { return fmt.Errorf("load config: %w", err) })

// Either: one of two types.
var e fp.Either[string, int] = fp.Right[string](42)
label := e.Fold(strconv.Quote, strconv.Itoa)
```

Prefer plain `(T, bool)` and `(T, error)` returns in public APIs; reach for these types where chaining makes the code clearer.

## Examples

Runnable programs live in [`examples/`](examples):

| Example | Shows |
| --- | --- |
| [`option`](examples/option/main.go) | lookups, chaining across types, JSON encoding |
| [`result`](examples/result/main.go) | chaining `(value, error)` functions, recovery, error context |
| [`either`](examples/either/main.go) | folding and mapping both sides |

```sh
go run ./examples/result
```

## Design

fp is inspired by [samber/mo](https://github.com/samber/mo) and Rust's `Option` and `Result`. It targets Go 1.27 only, which allows a smaller API: generic methods replace the `MapOption`/`MapResult`-style helpers that older Go versions need, and iterator and pointer interop come from the standard library.

## License

[MIT](LICENSE)
