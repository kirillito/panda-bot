package db

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	if e := (&Error{"Foo", errors.New("Bar")}); e.Error() != "Foo: Bar" {
		t.Errorf(`unexpected error %q`, e)
	}

	if e := (&Error{message: "Foo"}); e.Error() != "Foo" {
		t.Errorf(`unexpected error %q`, e)
	}
}