// Command option models values that may be absent with fp.Option, chains
// them with Go 1.27 generic methods, and encodes them as JSON.
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/luminyra/fp"
)

// Profile has an optional Nickname, encoded as null or omitted in JSON.
type Profile struct {
	Name     string            `json:"name"`
	Nickname fp.Option[string] `json:"nickname"`
}

func main() {
	env := map[string]string{"HOST": "  example.com ", "MODE": "   ", "WORKERS": "8"}

	// region: lookup
	// Turn the comma-ok idiom into an Option and treat blank values as absent.
	lookup := func(key string) fp.Option[string] {
		v, ok := env[key]
		return fp.OptionOf(strings.TrimSpace(v), ok).Filter(func(s string) bool { return s != "" })
	}

	host := lookup("HOST").Or("localhost")
	mode := lookup("MODE").Or("production")
	// endregion

	fmt.Println("host:", host)
	fmt.Println("mode:", mode)

	// region: chain
	// Map and AndThen are generic methods, so the chain turns a string into an
	// int directly.
	parseInt := func(s string) fp.Option[int] {
		n, err := strconv.Atoi(s)
		return fp.OptionOf(n, err == nil)
	}
	workers := lookup("WORKERS").AndThen(parseInt).Map(func(n int) int { return n * 2 })
	debug := lookup("DEBUG").Map(func(s string) bool { return s == "1" })
	// endregion

	fmt.Println("workers:", workers, "->", workers.Or(1))
	fmt.Println("debug:  ", debug, "->", debug.Or(false))

	// region: json
	// An empty Option encodes as null; null or a missing field decodes as None.
	out, err := json.Marshal([]Profile{
		{Name: "alice", Nickname: fp.Some("ali")},
		{Name: "bob"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))

	var p Profile
	if err := json.Unmarshal([]byte(`{"name":"carol"}`), &p); err != nil {
		panic(err)
	}
	fmt.Println(p.Name, p.Nickname, p.Nickname.IsNone())
	// endregion
}
