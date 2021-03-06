package ango

import (
	"encoding/json"
)

type messageIn struct {
	// Type of information contained in this message
	// See protocol.md for specification of values
	Type string `json:"type"`

	// Procedure name for the to-be-called procedure
	// exclusively for type "req"
	Procedure string `json:"procedure,omitempty"`

	// cb_id is a callback generated by side that creates the req
	// when set for type "req", the client expects a res/rej, and can handle not.
	// mandatory for types "res", "rej", "not"
	CallbackID uint64 `json:"cb_id"`

	// def_id is a reference to a deferred
	DeferredID uint64 `json:"def_id"`

	// lo_id refers to a registered linked object
	// it is set for types "lora" and "lou"
	LinkedObjectID uint64 `json:"lo_id"`

	// data object
	// type "req": parameters for the to-be-called procedure
	// types "rej", "res", "not": data for the resolve/reject/notify
	// type "lou": the new version of the linked object
	// optional for types "rej", "res", "not"
	Data json.RawMessage `json:"data"`

	// error string
	Error string `json:"error"`
}

type messageOut struct {
	// Type of information contained in this message
	// See protocol.md for specification of values
	Type string `json:"type"`

	// Procedure name for the to-be-called procedure
	// exclusively for type "req"
	Procedure string `json:"procedure,omitempty"`

	// cb_id is a callback generated by side that creates the req
	// when set for type "req", the client expects a res/rej, and can handle not.
	// mandatory for types "res", "rej", "not"
	CallbackID *uint64 `json:"cb_id,omitempty"`

	// def_id is a reference to a deferred
	DeferredID *uint64 `json:"def_id,omitempty"`

	// lo_id refers to a registered linked object
	// it is set for types "lora" and "lou"
	LinkedObjectID *uint64 `json:"lo_id,omitempty"`

	// data object
	// type "req": parameters for the to-be-called procedure
	// types "rej", "res", "not": data for the resolve/reject/notify
	// type "lou": the new version of the linked object
	// optional for types "rej", "res", "not"
	Data interface{} `json:"data,omtempty"`

	// error string
	Error string `json:"error,omitempty"`
}
