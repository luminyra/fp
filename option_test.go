package fp

import (
	"encoding/json"
	"slices"
	"strconv"
	"testing"
)

func TestOption(t *testing.T) {
	some, none := Some(2), None[int]()

	if v, ok := some.Get(); !ok || v != 2 || !some.IsSome() || some.IsNone() {
		t.Error("Some mismatch")
	}
	if _, ok := none.Get(); ok || !none.IsNone() {
		t.Error("None mismatch")
	}
	var zero Option[int]
	if zero.IsSome() {
		t.Error("zero Option must be empty")
	}
	if none.Or(7) != 7 || some.Or(7) != 2 || none.OrZero() != 0 {
		t.Error("Or/OrZero mismatch")
	}
	calls := 0
	fallback := func() int { calls++; return 9 }
	if some.OrFunc(fallback) != 2 || none.OrFunc(fallback) != 9 || calls != 1 {
		t.Error("OrFunc must call fallback only when empty")
	}
	if none.Pointer() != nil || *some.Pointer() != 2 {
		t.Error("Pointer mismatch")
	}
	if OptionFromPointer[int](nil).IsSome() || OptionFromPointer(some.Pointer()).MustGet() != 2 {
		t.Error("OptionFromPointer mismatch")
	}
	if OptionOf(1, false).IsSome() || !OptionOf(1, true).IsSome() {
		t.Error("OptionOf mismatch")
	}
	if some.Filter(func(v int) bool { return v > 5 }).IsSome() || !some.Filter(func(v int) bool { return v > 1 }).IsSome() {
		t.Error("Filter mismatch")
	}
	if some.String() != "Some(2)" || none.String() != "None" {
		t.Error("String mismatch")
	}
	if got := slices.Collect(some.All()); !slices.Equal(got, []int{2}) {
		t.Errorf("Some.All = %v", got)
	}
	if got := slices.Collect(none.All()); len(got) != 0 {
		t.Errorf("None.All = %v", got)
	}

	defer func() {
		if r := recover(); r != ErrNone {
			t.Fatalf("MustGet panic = %v", r)
		}
	}()
	none.MustGet()
}

func TestOptionMapChangesType(t *testing.T) {
	half := func(n int) Option[int] { return OptionOf(n/2, n%2 == 0) }

	tests := []struct {
		name     string
		input    Option[int]
		expected Option[string]
	}{
		{name: "even", input: Some(8), expected: Some("4")},
		{name: "odd", input: Some(7), expected: None[string]()},
		{name: "empty", input: None[int](), expected: None[string]()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.AndThen(half).Map(strconv.Itoa)
			if got != tt.expected {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestOptionJSON(t *testing.T) {
	type payload struct {
		A Option[int] `json:"a"`
		B Option[int] `json:"b"`
	}
	data, err := json.Marshal(payload{A: Some(1)})
	if err != nil || string(data) != `{"a":1,"b":null}` {
		t.Fatalf("Marshal = %s, %v", data, err)
	}
	var out payload
	if err := json.Unmarshal([]byte(`{"a":null,"b":3}`), &out); err != nil {
		t.Fatal(err)
	}
	if out.A.IsSome() || out.B.MustGet() != 3 {
		t.Fatalf("Unmarshal = %+v", out)
	}
	if err := json.Unmarshal([]byte(`{"a":"x"}`), &out); err == nil {
		t.Fatal("Unmarshal accepted a string for Option[int]")
	}
}
