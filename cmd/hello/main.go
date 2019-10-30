package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/marco-m/travis-go-dockerhub/hello"
)

func main() {
	var args struct {
		Foo string `default:"bar"`
	}
	arg.MustParse(&args)
	fmt.Println(hello.Hello())
}
