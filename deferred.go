package ango

import (
	"errors"
)

// Deferred provides methods to update a deferred in the angularjs client
type Deferred struct {
	conn       Conn
	deferredID uint64
	done       bool
}

// ErrAlreadyResolvedOrRejected is returned by Resolve() or Reject() when the Deferred on which the method is called has been resolved or rejected already
var ErrAlreadyResolvedOrRejected = errors.New("deferred was already resolved or rejected")

// Resolve the deferred in the AngularJS client
func (d *Deferred) Resolve(data interface{}) error {
	// check if deferred isn't already completed
	if d.done {
		return ErrAlreadyResolvedOrRejected
	}
	d.done = true

	// send message with type res and def_id
	err := d.conn.sendMessage(&messageOut{
		Type:       "res",
		DeferredID: &d.deferredID,
		Data:       data,
	})
	if err != nil {
		return err
	}

	// all done
	return nil
}

// Reject the deferred in the AngularJS client
func (d *Deferred) Reject(data interface{}) error {
	// check if deferred isn't already completed
	if d.done {
		return ErrAlreadyResolvedOrRejected
	}
	d.done = true

	// send message with type rej and def_id
	err := d.conn.sendMessage(&messageOut{
		Type:       "rej",
		DeferredID: &d.deferredID,
		Data:       data,
	})
	if err != nil {
		return err
	}

	// all done
	return nil
}

// Notify the deferred in the AngularJS client
func (d *Deferred) Notify(data interface{}) error {
	// send message with type not and def_id
	err := d.conn.sendMessage(&messageOut{
		Type:       "not",
		DeferredID: &d.deferredID,
		Data:       data,
	})
	if err != nil {
		return err
	}

	// all done
	return nil
}
