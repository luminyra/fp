// Command result chains ordinary (value, error) functions with fp.Result and
// recovers from or annotates a failure.
package main

import (
	"fmt"
	"strconv"

	"github.com/luminyra/fp"
)

// checkPort is an ordinary Go function that knows nothing about fp.
func checkPort(n int) (int, error) {
	if n < 1 || n > 65535 {
		return 0, fmt.Errorf("port %d out of range", n)
	}
	return n, nil
}

// region: chain
// parsePort parses an integer, then checks its range; when a step fails, the
// following steps are skipped. AndThen takes a (value, error) function as is.
func parsePort(s string) fp.Result[int] {
	return fp.ResultOf(strconv.Atoi(s)).AndThen(checkPort)
}

// endregion

func main() {
	for _, s := range []string{"8080", "http", "70000"} {
		fmt.Printf("%-6s -> %v\n", s, parsePort(s))
	}

	// region: recover
	// Recover runs only on failure and can supply a fallback; MapErr adds
	// context to the error.
	port := parsePort("http").Recover(func(error) (int, error) { return 8080, nil })
	fmt.Println("fallback port:", port.MustGet())

	_, err := parsePort("70000").MapErr(func(err error) error {
		return fmt.Errorf("read PORT: %w", err)
	}).Get()
	fmt.Println("error:", err)
	// endregion

	// region: convert
	// Map is a generic method, so a Result[int] becomes a Result[string]; Get
	// returns to the plain (value, error) form.
	addr, err := parsePort("443").Map(func(n int) string { return ":" + strconv.Itoa(n) }).Get()
	fmt.Println("listen address:", addr, err)
	// endregion
}
