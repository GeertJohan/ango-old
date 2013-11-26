package main

import (
	"github.com/GeertJohan/ango"
)

type SomeStuff struct {
	ango.LinkedObject // not very nice that it must be a pointer
	Foo               string
	Bar               string
}

func main() {
	ss := SomeStuff{
		Foo: "foo",
		Bar: "bar",
	}
	ango.LinkedObjectEater(ss)
}
