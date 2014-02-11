## Ango

**In development**, look arround but don't use this just yet :)

**NOTE: Development on this project has stopped. This project has been renamed and is marked for deletion. A new project also named 'ango' has emerged, with essentially the same goal as this one, but a totally different approach. For more information, visit: https://github.com/GeertJohan/ango**

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
