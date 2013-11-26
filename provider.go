package ango

import (
	"code.google.com/p/go.net/websocket"
	"github.com/GeertJohan/go.incremental"
	"net/http"
)

// Provider implements the http.Handler interface
type Provider struct {
	wsHandler websocket.Handler

	// conn id counter
	idCounter incremental.Uint64

	// BeforeWebsocketCheck can be set to perform additional checks on a new connection (e.g. auth, check origin, etc.)
	// When false is returned, the connection is dropped. The funciton itself is responsible for correcly replying to the client with a http status code
	BeforeWebsocketCheck func(w http.ResponseWriter, r *http.Request) bool
}

// NewProvider create a new Provider instance
func NewProvider() *Provider {
	p := &Provider{}
	p.wsHandler = websocket.Handler(p.setupWebsocket)
	return p
}

func (p *Provider) setupWebsocket(wsConn *websocket.Conn) {
	c := &Conn{
		provider: p,
		connID:   p.idCounter.Next(),
		conn:     wsConn,
	}

	c.receiveAndHandle()
}

// ServeHTTP implements the http.Handler interface
// It provides a websocket for each incomming request
func (p *Provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// perform BeforeWebsocketCheck
	if p.BeforeWebsocketCheck != nil && !p.BeforeWebsocketCheck(w, r) {
		return
	}
	// start websocket
	p.wsHandler.ServeHTTP(w, r)
}
