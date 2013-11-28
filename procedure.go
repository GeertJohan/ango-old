package ango

import (
	"encoding/json"
)

// ProcedureFunc can be implemented to handle
type ProcedureFunc func(data json.RawMessage, def *Deferred)

// RegisterProcedureFunc saves a named procedure to the given provider, which called from that point on.
func (p *Provider) RegisterProcedureFunc(name string, fn ProcedureFunc) {
	p.procedures[name] = fn
}
