package fp

import (
	"strconv"
	"testing"
)

func TestEither(t *testing.T) {
	l := Left[string, int]("x")
	r := Right[string](3)

	if !l.IsLeft() || l.IsRight() || !r.IsRight() {
		t.Error("side mismatch")
	}
	if l.Left() != Some("x") || l.Right().IsSome() || r.Right() != Some(3) || r.Left().IsSome() {
		t.Error("Left/Right mismatch")
	}
	if s := l.Swap(); s.Right() != Some("x") {
		t.Error("Swap mismatch")
	}
	if got := r.Fold(func(s string) string { return s }, strconv.Itoa); got != "3" {
		t.Errorf("Fold = %q", got)
	}
	if got := l.Fold(func(s string) int { return len(s) }, func(n int) int { return n }); got != 1 {
		t.Errorf("Fold = %d", got)
	}
	if got := r.MapRight(strconv.Itoa); got.Right() != Some("3") {
		t.Errorf("MapRight = %v", got)
	}
	if got := l.MapRight(strconv.Itoa); got.Left() != Some("x") {
		t.Errorf("MapRight on left = %v", got)
	}
	if got := l.MapLeft(func(s string) int { return len(s) }); got.Left() != Some(1) {
		t.Errorf("MapLeft = %v", got)
	}
	if got := r.MapLeft(func(s string) int { return len(s) }); got.Right() != Some(3) {
		t.Errorf("MapLeft on right = %v", got)
	}
	if l.String() != "Left(x)" || r.String() != "Right(3)" {
		t.Error("String mismatch")
	}
}
