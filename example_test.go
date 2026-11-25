package fp_test

import (
	"fmt"
	"strconv"

	"github.com/luminyra/fp"
)

func ExampleOptionOf() {
	ports := map[string]string{"http": "80"}

	port, ok := ports["https"]
	fmt.Println(fp.OptionOf(port, ok).Or("443"))
	// Output: 443
}

func ExampleOption_AndThen() {
	parse := func(s string) fp.Option[int] {
		n, err := strconv.Atoi(s)
		return fp.OptionOf(n, err == nil)
	}

	fmt.Println(fp.Some("42").AndThen(parse).Map(func(n int) int { return n * 2 }))
	fmt.Println(fp.Some("x").AndThen(parse))
	// Output:
	// Some(84)
	// None
}

func ExampleOption_All() {
	for v := range fp.Some("only value").All() {
		fmt.Println(v)
	}
	for v := range fp.None[string]().All() {
		fmt.Println("never printed:", v)
	}
	// Output: only value
}

func ExampleResult_AndThen() {
	checkPort := func(n int) (int, error) {
		if n < 1 || n > 65535 {
			return 0, fmt.Errorf("port %d out of range", n)
		}
		return n, nil
	}

	for _, raw := range []string{"8080", "http", "70000"} {
		fmt.Println(fp.ResultOf(strconv.Atoi(raw)).AndThen(checkPort))
	}
	// Output:
	// Ok(8080)
	// Err(strconv.Atoi: parsing "http": invalid syntax)
	// Err(port 70000 out of range)
}

func ExampleResult_Recover() {
	port := fp.ResultOf(strconv.Atoi("http")).Recover(func(error) (int, error) { return 8080, nil })
	fmt.Println(port)
	// Output: Ok(8080)
}

func ExampleEither_Fold() {
	describe := func(e fp.Either[string, int]) string {
		return e.Fold(
			func(text string) string { return "text " + strconv.Quote(text) },
			func(n int) string { return "number " + strconv.Itoa(n) },
		)
	}

	fmt.Println(describe(fp.Left[string, int]("hello")))
	fmt.Println(describe(fp.Right[string](42)))
	// Output:
	// text "hello"
	// number 42
}
