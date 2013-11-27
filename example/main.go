package main

import (
	"github.com/GeertJohan/ango"
	"log"
	"net/http"
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
	c := &ango.Conn{}
	c.Call("hi", nil, nil, nil, nil)

	http.Handle("/", http.FileServer(http.Dir("./files")))
	http.Handle("/ango.js", http.FileServer(http.Dir("../")))

	// create new provider, host it on
	p := ango.NewProvider()
	p.BeforeWebsocket = func(w http.ResponseWriter, r *http.Request) bool {
		log.Printf("new incomming connection!")
		return true
	}
	http.Handle("/ango-websocket", p)

	// listen and serve
	http.ListenAndServe(":8123", nil)
}
