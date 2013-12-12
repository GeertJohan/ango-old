package ango

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GeertJohan/go.incremental"
	"io"
	"log"
	"sync"
)

// Conn abstracts the websocket an provides methods to communicate with the client
type Conn struct {
	provider *Provider
	connID   uint64
	conn     *websocket.Conn

	// request callbacks
	callbackInc          incremental.Uint64
	callbackChannels     map[uint64](chan callback)
	callbackChannelsLock sync.Mutex

	// call promises
	promiseInc   incremental.Uint64
	promises     map[uint64]*promise
	promisesLock sync.Mutex
}

type callback struct {
	accepted bool   // accepted=true, denied=false
	err      string // optionally error, why it was denied
	newID    uint64 // client-side created ID, if any
}

// Fire is a fire-and-forget method to run a named procedure on the connected client
// The procedure result is never observed by this server. An error is returned when the procedure could not be called.
func (c *Conn) Fire(name string, data interface{}) error {
	// send message with type res and def_id
	err := c.sendRequest(&messageOut{
		Type:      "req",
		Procedure: name,
		Data:      data,
	})
	if err != nil {
		return err
	}

	// all done
	return nil
}

// PromiseFunc is a handlerfunc, called when a given promise is resolved/notified/rejected
type PromiseFunc func(json.RawMessage)

type promise struct {
	id        uint64
	resolveFn PromiseFunc
	rejectFn  PromiseFunc
	notifyFn  PromiseFunc
	waitCh    chan bool
}

func (c *Conn) newPromise(resolve PromiseFunc, reject PromiseFunc, notify PromiseFunc) *promise {
	p := &promise{
		id:        c.promiseInc.Next(),
		resolveFn: resolve,
		rejectFn:  reject,
		notifyFn:  notify,
		waitCh:    make(chan bool, 1),
	}

	c.promisesLock.Lock()
	c.promises[p.id] = p
	c.promisesLock.Unlock()

	return p
}

func (c *Conn) newCallbackChan() (uint64, chan callback) {
	// setup callback for the request (callback reqa/reqd and lora/lord)
	c.callbackChannelsLock.Lock()
	cbid := c.callbackInc.Next()
	cbch := make(chan callback, 1)
	c.callbackChannels[cbid] = cbch
	c.callbackChannelsLock.Unlock()
	return cbid, cbch

}

func (p *promise) waitForFeedback() {
	<-p.waitCh
}

// CallAndWait runs the procedure and returns when the call is completed and resolve/reject was ran
// An error is returned immediatly when the procedure with given name could not be called (name not registered in client or connection broken).
func (c *Conn) CallAndWait(name string, data interface{}, resolve PromiseFunc, reject PromiseFunc, notify PromiseFunc) error {
	return c.call(true, name, data, resolve, reject, notify)
}

// Call runs the procedure and returns when the request for a procedure call has been made, but has not necicarily completed yet.
// The resolve/reject/notify function can be called later-on. An error is returned when the given procedure name does not exist.
// An error is returned immediatly when the procedure with given name could not be called (name not registered in client or connection broken).
func (c *Conn) Call(name string, data interface{}, resolve PromiseFunc, reject PromiseFunc, notify PromiseFunc) error {
	return c.call(false, name, data, resolve, reject, notify)
}

func (c *Conn) call(sync bool, name string, data interface{}, resolve PromiseFunc, reject PromiseFunc, notify PromiseFunc) error {
	if resolve == nil || reject == nil {
		return errors.New("resolve and reject must be set with a valid PromiseFunc")
	}

	// create a new promise
	p := c.newPromise(resolve, reject, notify)

	// send message with type res and def_id
	err := c.sendRequest(&messageOut{
		Type:       "req",
		Procedure:  name,
		Data:       data,
		DeferredID: &p.id,
	})
	if err != nil {
		return err
	}

	// sync or async
	if sync {
		p.waitForFeedback()
	} else {
		go p.waitForFeedback()
	}

	// all done
	return nil
}

func (c *Conn) sendRequest(msg *messageOut) error {
	// create new callback and link to message
	cbid, cbch := c.newCallbackChan()
	msg.CallbackID = &cbid

	// send message
	err := c.sendMessage(msg)
	if err != nil {
		return err
	}

	//++ TODO: add timeout?
	cb := <-cbch
	if cb.accepted {
		return nil
	}
	return fmt.Errorf("request denied, %s", cb.err)
}

func (c *Conn) sendMessage(msg *messageOut) error {

	// try data to be a linked object
	if msg.Data != nil {
		loid, err := c.tryLinkedObject(msg.Data)
		if err != nil {
			return err
		}
		if loid != nil {
			msg.Data = nil
			msg.LinkedObjectID = loid
		}
	}

	// send message over websocket
	err := websocket.JSON.Send(c.conn, msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) sendRequestResponse(typ string, cbid uint64) error {
	out := &messageOut{
		Type:       typ,
		CallbackID: &cbid,
	}
	return c.sendMessage(out)
}

func (c *Conn) receiveAndHandle() {
	// loop and receive messages
	for {
		in := &messageIn{}
		err := websocket.JSON.Receive(c.conn, in)
		if err != nil {
			if c.provider.Debug {
				log.Printf("Closing connection %d.\n", c.connID)
			}
			//++ clean up

			// panic if error isn't os.EOF
			if err != io.EOF {
				panic(err)
			}
			return // closes connection
		}
		switch in.Type {
		case "req":
			c.provider.proceduresLock.RLock()
			proc, ok := c.provider.procedures[in.Procedure]
			c.provider.proceduresLock.RUnlock()
			if !ok {
				if c.provider.Debug {
					log.Printf("Request for procedure '%s', but that is not registered.\n", in.Procedure)
				}
				c.sendRequestResponse("reqd", in.CallbackID)
				continue
			}
			c.sendRequestResponse("reqa", in.CallbackID)
			deferred := &Deferred{
				conn:       c,
				deferredID: in.DeferredID,
			}
			proc(in.Data, deferred)
			if !deferred.done {
				deferred.Reject("procedure did not resolve nor reject")
			}

		case "reqa", "reqd", "lora", "lord":
			// find callback channel
			c.callbackChannelsLock.Lock()
			cbch, ok := c.callbackChannels[in.CallbackID]
			if !ok {
				if c.provider.Debug {
					fmt.Printf("Got msg w/ type: \"%s\", cb_id: %d but could not find handler\n", in.Type, in.CallbackID)
				}
				c.callbackChannelsLock.Unlock()
				continue
			}
			delete(c.callbackChannels, in.CallbackID)
			c.callbackChannelsLock.Unlock()

			switch in.Type {
			case "reqa":
				cbch <- callback{
					accepted: true,
				}
			case "reqd":
				cbch <- callback{
					accepted: false,
					err:      in.Error,
				}
			case "lora":
				cbch <- callback{
					accepted: true,
					newID:    in.LinkedObjectID,
				}
			case "lord":
				//++ TODO: this never really happens, now does it?
				cbch <- callback{
					accepted: false,
					err:      in.Error,
				}
			}

		case "res", "rej":
			c.promisesLock.Lock()
			p, ok := c.promises[in.DeferredID]
			if !ok {
				if c.provider.Debug {
					fmt.Printf("Got msg w/ type: \"%s\", def_id: %d but could not find promise\n", in.Type, in.DeferredID)
				}
				c.callbackChannelsLock.Unlock()
				continue
			}
			delete(c.promises, in.DeferredID)
			c.promisesLock.Unlock()

			// fullfill based on type
			switch in.Type {
			case "res":
				p.resolveFn(in.Data)
			case "rej":
				p.rejectFn(in.Data)
			case "not":
				p.notifyFn(in.Data)
			}

			if in.Type != "not" {
				// the wait for res/rej is over
				p.waitCh <- true
			}

		default:
			if c.provider.Debug {
				fmt.Printf("Got msg w/ unknown type: \"%s\"\n", in.Type)
			}
			continue
		}
	}

	// when this returns, connection is closed
}
