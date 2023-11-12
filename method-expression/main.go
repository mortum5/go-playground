package main

import (
	"fmt"
	"log"
)

type T struct{}

func (T) Foo(s string) { println(s) }
func (T) Bar(s string) { log.Println(s) }
func (T) Baz(s string) { fmt.Printf("%s", s) }

type F func(T, string)

func main() {
	t := T{}
	for _, fn := range []F{T.Foo, T.Bar, T.Baz} {
		fn(t, "Hello")
	}
}
