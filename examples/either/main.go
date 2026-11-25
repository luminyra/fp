// Command either holds one of two value types with fp.Either and transforms
// and folds both sides with generic methods.
package main

import (
	"fmt"
	"strconv"

	"github.com/luminyra/fp"
)

func main() {
	// region: either
	// Input that parses as a number goes right; anything else stays left.
	parse := func(s string) fp.Either[string, int] {
		if n, err := strconv.Atoi(s); err == nil {
			return fp.Right[string](n)
		}
		return fp.Left[string, int](s)
	}

	for _, s := range []string{"42", "hello"} {
		e := parse(s)
		desc := e.Fold(
			func(text string) string { return "text " + strconv.Quote(text) },
			func(n int) string { return fmt.Sprintf("number %d, doubled %d", n, n*2) },
		)
		fmt.Printf("%-13v -> %s\n", e, desc)
	}
	// endregion

	// region: map
	// MapRight transforms only the right side and may change its type; Right
	// returns the right value as an Option.
	doubled := parse("21").MapRight(func(n int) int { return n * 2 })
	fmt.Println(doubled, doubled.Right())
	fmt.Println(parse("7").Swap())
	// endregion
}
