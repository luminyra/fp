package fp

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestResult(t *testing.T) {
	boom := errors.New("boom")
	ok, bad := Ok(2), Err[int](boom)

	if !ok.IsOk() || ok.IsErr() || ok.Err() != nil {
		t.Error("Ok mismatch")
	}
	if !bad.IsErr() || !errors.Is(bad.Err(), boom) {
		t.Error("Err mismatch")
	}
	if Err[int](nil).IsOk() {
		t.Error("Err(nil) must still be a failure")
	}
	if v, err := ok.Get(); v != 2 || err != nil {
		t.Error("Get mismatch")
	}
	if bad.Or(5) != 5 || bad.OrZero() != 0 || ok.OrZero() != 2 {
		t.Error("Or/OrZero mismatch")
	}
	if ResultOf(1, boom).IsOk() || !ResultOf(strconv.Atoi("1")).IsOk() {
		t.Error("ResultOf mismatch")
	}
	if !errors.Is(Errf[int]("wrap: %w", boom).Err(), boom) {
		t.Error("Errf must keep the error chain")
	}
	if ok.Option().MustGet() != 2 || bad.Option().IsSome() {
		t.Error("Option mismatch")
	}
	if ok.String() != "Ok(2)" || bad.String() != "Err(boom)" {
		t.Error("String mismatch")
	}

	defer func() {
		if r := recover(); r != boom {
			t.Fatalf("MustGet panic = %v", r)
		}
	}()
	bad.MustGet()
}

func TestResultChain(t *testing.T) {
	boom := errors.New("boom")
	checkPort := func(n int) (int, error) {
		if n < 1 || n > 65535 {
			return 0, fmt.Errorf("port %d out of range", n)
		}
		return n, nil
	}

	addr := ResultOf(strconv.Atoi("443")).AndThen(checkPort).Map(func(n int) string { return ":" + strconv.Itoa(n) })
	if addr.MustGet() != ":443" {
		t.Errorf("chain = %v", addr)
	}

	called := false
	failed := Err[int](boom).AndThen(func(int) (int, error) { called = true; return 0, nil })
	if called || !errors.Is(failed.Err(), boom) {
		t.Error("AndThen must skip f and keep the error on failure")
	}
	if !errors.Is(Err[int](boom).Map(strconv.Itoa).Err(), boom) {
		t.Error("Map must keep the error on failure")
	}

	wrapped := Err[int](boom).MapErr(func(err error) error { return fmt.Errorf("load: %w", err) })
	if wrapped.Err().Error() != "load: boom" || !errors.Is(wrapped.Err(), boom) {
		t.Errorf("MapErr = %v", wrapped)
	}
	if Ok(1).MapErr(func(error) error { return boom }).IsErr() {
		t.Error("MapErr must leave a success unchanged")
	}
	if Err[int](boom).MapErr(func(error) error { return nil }).IsOk() {
		t.Error("MapErr returning nil must still be a failure")
	}

	recovered := Err[int](boom).Recover(func(error) (int, error) { return 9, nil })
	if recovered.MustGet() != 9 || Ok(2).Recover(func(error) (int, error) { return 0, boom }).MustGet() != 2 {
		t.Error("Recover mismatch")
	}
}
