package hello_test

import (
	"testing"

	"github.com/marco-m/travis-go-dockerhub/hello"
)

func TestHello(t *testing.T) {
	want := "42"
	if got := hello.Hello(); got != want {
		t.Fatalf("got %q; want %q", got, want)
	}
}
