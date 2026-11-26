<img align="right" src="https://images.luminyra.dev/luminyra/github/fp/logo.webp" alt="A sage path with ascending folds and a short terracotta failure path." width="96" height="96">

# fp

[![Go Reference](https://pkg.go.dev/badge/github.com/luminyra/fp.svg)](https://pkg.go.dev/github.com/luminyra/fp)
[![CI](https://github.com/luminyra/fp/actions/workflows/ci.yml/badge.svg)](https://github.com/luminyra/fp/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/luminyra/fp)](https://goreportcard.com/report/github.com/luminyra/fp)
[![Go 1.27+](https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/doc/devel/release)

fp gives Go three small types for values that are not simply "there": **`Option`** for a value that may be missing, **`Result`** for a value that may have failed, and **`Either`** for a value of one of two types.

They let you write "parse this, then check that, then format it, and fall back to a default if anything goes wrong" as one readable chain, while still working with the ordinary `(value, error)` and `(value, ok)` functions you already have.

- **Chains that change type.** `Map` and `AndThen` are Go 1.27 generic methods, so a `Result[int]` can become a `Result[string]` mid-chain.
- **Plain Go in, plain Go out.** `ResultOf(strconv.Atoi(s))` and `OptionOf(value, ok)` bring existing code in; `Get` turns a chain back into `(value, error)` or `(value, ok)`.
- **Fits the standard library.** `Option` encodes to JSON as its value or `null` and can be ranged over as an `iter.Seq`.
- **Small and dependency-free.** Three types, immutable values, no third-party code.
- **No overhead.** A `Result` chain runs as fast as the same steps written with `if err != nil`.

<p align="center">
  <img src="https://images.luminyra.dev/luminyra/github/fp/cover.webp" alt="A sage-green mechanical track connects several stages above a coral path." width="640">
</p>

## Contents

- [When to use it](#when-to-use-it) · [Install](#install) · [Getting started](#getting-started)
- [Common tasks](#common-tasks) · [How it works](#how-it-works) · [Performance](#performance) · [API at a glance](#api-at-a-glance)

## When to use it

fp helps where "missing" and "failed" are a normal part of the flow:

- reading configuration with defaults and validation,
- parsing and converting input in several steps,
- optional fields in JSON documents,
- code that tries several sources and takes the first that works.

It is not a replacement for Go's error handling. In public APIs, return plain `(T, error)` or `(T, bool)`; a single `if err != nil` is clearer than a chain of one step.

## Install

```sh
go get github.com/luminyra/fp
```

fp requires Go 1.27 or later and has no dependencies.

## Getting started

The program reads settings from a map, as it might from the environment. Copy it into `main.go` and run it with `go run .`:

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/luminyra/fp"
)

// checkPort is an ordinary Go function; fp does not need anything special.
func checkPort(n int) (int, error) {
	if n < 1 || n > 65535 {
		return 0, fmt.Errorf("port %d out of range", n)
	}
	return n, nil
}

func main() {
	env := map[string]string{"PORT": "8080", "WORKERS": "many"}

	// Option: a value that may be missing.
	lookup := func(key string) fp.Option[string] {
		value, ok := env[key]
		return fp.OptionOf(value, ok)
	}
	host := lookup("HOST").Or("localhost")

	// Result: a value that may have failed. Each step runs only if the
	// previous one succeeded, and Map can change the type along the way.
	addr := fp.ResultOf(strconv.Atoi(lookup("PORT").Or("80"))).
		AndThen(checkPort).
		Map(func(port int) string { return host + ":" + strconv.Itoa(port) })

	workers := fp.ResultOf(strconv.Atoi(lookup("WORKERS").Or("4")))

	fmt.Println("address:", addr)
	fmt.Println("workers:", workers)
	fmt.Println("workers, with a fallback:", workers.Or(4))
}
```

```text
address: Ok(localhost:8080)
workers: Err(strconv.Atoi: parsing "many": invalid syntax)
workers, with a fallback: 4
```

`ResultOf` turned `strconv.Atoi`'s `(int, error)` into a `Result`, `AndThen` ran `checkPort` only because parsing succeeded, and `Map` turned the port into a string. The failed `WORKERS` value kept its original error until `Or` replaced it with a default.

<p align="center">
  <a href="https://images.luminyra.dev/luminyra/github/fp/demo.webp"><img src="https://images.luminyra.dev/luminyra/github/fp/demo.webp" alt="A valid port passes through parsing, validation and formatting. A parse error skips the later steps." width="760"></a>
</p>

## Common tasks

### Turn a comma-ok lookup into an Option

```go
value, ok := settings["timeout"]
timeout := fp.OptionOf(value, ok).Or("30s")
```

`OptionFromPointer(p)` does the same for a pointer that may be nil, and `Option.Pointer` converts back.

### Compute a default only when it is needed

```go
name := fp.OptionOf(value, ok).OrFunc(defaultName) // defaultName runs only if the value is missing
```

### Chain steps that may fail

`AndThen` accepts any `func(T) (R, error)`, and `Map` any `func(T) R`:

```go
config := fp.ResultOf(os.ReadFile(path)).
	AndThen(parseConfig). // func([]byte) (Config, error)
	AndThen(validate)     // func(Config) (Config, error)
```

`Option` has the same pair: `AndThen` for steps that may produce nothing, and `Map` for steps that always do.

### Add context to an error, or recover from it

```go
config = config.MapErr(func(err error) error { return fmt.Errorf("load %s: %w", path, err) })

fallback := config.Recover(func(error) (Config, error) { return defaultConfig, nil })
```

`errors.Is` and `errors.As` keep working on the wrapped error.

### Go back to plain Go

```go
cfg, err := config.Get()    // (Config, error), like any Go function
user, ok := maybeUser.Get() // (User, bool)
```

`Result.Option` drops the error and returns an `Option`. `MustGet` panics instead of returning an error, which suits tests and `main`.

### Make JSON fields optional

```go
type Profile struct {
	Name     string            `json:"name"`
	Nickname fp.Option[string] `json:"nickname"` // null or missing decodes as None
}
```

### Hold one of two types

```go
var reply fp.Either[string, int] = fp.Right[string](42)

text := reply.Fold(strconv.Quote, strconv.Itoa) // one function per side
```

`MapLeft` and `MapRight` transform one side, and `Left` and `Right` return the side as an `Option`.

## How it works

- All three types are plain structs. They are immutable: methods return new values and never modify the receiver, so values are safe to share between goroutines as long as what they hold is.
- The zero `Option` is empty, the zero `Result` is `Ok` with the zero value of `T`, and the zero `Either` is a `Left` with the zero value of `L`.
- `Err(nil)` still produces a failed `Result`, so a forgotten error check cannot turn a failure into success.
- `MustGet` panics with `fp.ErrNone` on an empty `Option` and with the held error on a failed `Result`.

## Performance

Measured with the benchmarks in this repository (`make bench`; Go 1.27.4, linux/arm64, 5 CPUs, median of 6 runs). Run them on your hardware for numbers that matter to you.

| Code | Time | Allocations |
| --- | --- | --- |
| `fp.ResultOf(strconv.Atoi(s)).AndThen(checkPort).Map(format).Get()` | 24 ns | 2 |
| The same steps in plain Go | 24 ns | 2 |
| `fp.OptionOf(v, ok).AndThen(parse).Map(double).Or(1)` | 5.9 ns | 0 |

The types are small structs passed by value, so a chain costs what its steps cost: both versions above allocate only inside `strconv`.

## API at a glance

| Type | Create | Read | Transform |
| --- | --- | --- | --- |
| `Option[T]` | `Some`, `None`, `OptionOf`, `OptionFromPointer` | `Get`, `Or`, `OrFunc`, `OrZero`, `MustGet`, `IsSome`, `IsNone`, `Pointer`, `All` | `Map`, `AndThen`, `Filter` |
| `Result[T]` | `Ok`, `Err`, `Errf`, `ResultOf` | `Get`, `Err`, `Or`, `OrZero`, `MustGet`, `IsOk`, `IsErr`, `Option` | `Map`, `AndThen`, `MapErr`, `Recover` |
| `Either[L, R]` | `Left`, `Right` | `Left`, `Right`, `IsLeft`, `IsRight` | `Fold`, `MapLeft`, `MapRight`, `Swap` |

The full reference with runnable examples is on [pkg.go.dev](https://pkg.go.dev/github.com/luminyra/fp). Complete programs live in [`examples/`](examples): [`option`](examples/option/main.go), [`result`](examples/result/main.go) and [`either`](examples/either/main.go).

## Versioning

fp follows [semantic versioning](https://semver.org). Until v1.0.0, a minor release may change the API; every change is listed in the [changelog](CHANGELOG.md). It supports the two most recent Go releases, starting with Go 1.27.

## Contributing

Questions, bug reports and pull requests are welcome. See the [contributing guide](https://github.com/luminyra/.github/blob/main/CONTRIBUTING.md), and report security issues as described in the [security policy](https://github.com/luminyra/.github/blob/main/SECURITY.md).

fp was inspired by [samber/mo](https://github.com/samber/mo) and Rust's `Option` and `Result`. It is released under the [MIT License](LICENSE).
