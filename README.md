## Ango

**In development**, look arround but don't use this just yet :)

### RPC

### Linked Object (Server > Client sync)
 - links with multiple clients
 - only synchronized when explicitly asked (`ango.Sync(yourObject)`)
 - complete object is re-sent
 - object should not be modified by the client

### TwoWay Object (Server > Client and Client > Server sync)
 - one server, one client
 - complete object is re-sent on every update
 - write-lock

[![GoDoc](http://godoc.org/github.com/GeertJohan/ango?status.png)](http://godoc.org/github.com/GeertJohan/ango)