package main

import (
	"encoding/json"
	"github.com/GeertJohan/ango"
	"log"
	"net/http"
)

type SomeStuff struct {
	ango.LinkedObject
	Foo string
	Bar string
}

func main() {

	// serve static files and the latest ango.js
	http.Handle("/", http.FileServer(http.Dir("./files")))
	http.Handle("/ango.js", http.FileServer(http.Dir("../")))

	// create and setup new provider
	p := ango.NewProvider()
	p.BeforeWebsocket = func(w http.ResponseWriter, r *http.Request) bool {
		log.Printf("new incomming connection!")
		return true
	}
	p.Debug = true

	p.RegisterProcedureFunc("getStuff", func(data json.RawMessage, def *ango.Deferred) {
		// create SomeStuff
		stuff := SomeStuff{
			Foo: "foo",
			Bar: "bar",
		}

		// resolve defered with stuff
		def.Resolve(stuff)
	})

	// hook provider onto a custom url
	http.Handle("/ango-websocket", p)

	// listen and serve
	http.ListenAndServe(":8123", nil)
}
