package ango

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"errors"
	"github.com/GeertJohan/go.incremental"
)

// Conn abstracts the websocket an provides methods to communicate with the client
type Conn struct {
	provider  *Provider
	connID    uint64
	conn      *websocket.Conn
	cbCounter incremental.Uint64
	callbacks map[int]func(data json.RawMessage)
}

func (c *Conn) registerLinkedObject(lo interface{}) (int64, error) {
	//++ check that lo is actually a linkedObject
	return 0, errors.New("not implemented yet")
}

// Call is a fire-and-forget method to call a named service on the connected client
// The service result is never observed by this server. An error is returned when the service could not be called.
func (c *Conn) Call(name string, data interface{}) error {
	//++
	return errors.New("not implemented yet")
}

// PromiseFunc is a handlerfunc, called when a given promise is resolved/notified/rejected
type PromiseFunc func(json.RawMessage)

// Request calls the service and blocks until a response is returned
// An error is returned immediatly when the service with given name could not be called (name not registered in client or connection broken).
func (c *Conn) Request(name string, data interface{}, resolve PromiseFunc, notify PromiseFunc, reject PromiseFunc) error {
	return errors.New("not implemented yet")
}

func (c *Conn) sendMessage(msg *messageOut) error {
	err := websocket.JSON.Send(c.conn, msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) receiveAndHandle() {
	for {
		in := &messageIn{}
		err := websocket.JSON.Receive(c.conn, in)
		if err != nil {
			panic(err)
		}
		switch in.Type {
		// case "lor": // not for client>server yet

		case "lora":
			//++ finish linked object registration

		// case "lou": // not for client>server yet

		case "req":
			//++ we have a request!
			//++ switch on header, send reqd if not exists
			//++ create a Deferred and send reqa
			//++ call service with given details and the Deferred
			// all done

		case "reqa":
			//++ a request was accepted!

		case "reqd":
			//++ a request was denied!

		case "res":
			//++ we have a resolve

		case "rej":
			//++ we have a rejection

		case "not":
			//++ we have a notification
		default:
			// unknown message type, protocol error, will now close connection
			return
		}
	}

	// when this returns, connection is closed
}
