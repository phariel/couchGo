package couchGo

//Private

func init() {
	srv = &CouchServer{ip: "127.0.0.1", port: 5984, defaultDbName: "master"}
}

//Public
var srv *CouchServer

//CouchServer --
type CouchServer struct {
	ip            string
	port          int
	defaultDbName string
}

//ConfigServer -- input couch server configuration manually
func ConfigServer(cs *CouchServer) {
	srv = cs
}

//Update -- equals "Insert" without Id passing
func Update(rBatch []RequestData, cb func(res []ResponseData)) {
	batchRequests("PUT", rBatch, cb)
}

func Get(rBatch []RequestData, cb func(res []ResponseData)) {
	batchRequests("GET", rBatch, cb)
}
