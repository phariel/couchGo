package couch

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/satori/go.uuid"
)

var srv *CouchServer

type CouchServer struct {
	ip            string
	port          int
	defaultDbName string
}

type couchRequest struct {
	method   string
	dbName   string
	uuid     uuid.UUID
	jsonBody []byte
}

type CouchResponse struct {
	Err    bool
	Status int
	Body   string
}

func init() {
	srv = &CouchServer{ip: "127.0.0.1", port: 5984, defaultDbName: "master"}
}

func ConfigServer(cs *CouchServer) {
	srv = cs
}

func (c *couchRequest) doRequest() *CouchResponse {
	var res *CouchResponse
	if c.dbName != "" {
		urlArray := []string{"http://", srv.ip, ":", strconv.Itoa(srv.port), "/", c.dbName, "/", c.uuid.String()}
		url := strings.Join(urlArray, "")
		req, _ := http.NewRequest(c.method, url, bytes.NewBuffer(c.jsonBody))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			res = &CouchResponse{Err: false, Status: resp.StatusCode, Body: buf.String()}
		}
	} else {
		res = &CouchResponse{Err: true, Status: 403, Body: "{\"error\":true}"}
	}
	return res
}

func Insert(json []byte, dbName string) *CouchResponse {
	var res *CouchResponse
	if dbName == "" {
		dbName = srv.defaultDbName
	}

	cr := &couchRequest{method: "PUT", dbName: dbName, uuid: uuid.NewV4(), jsonBody: json}

	res = cr.doRequest()
	return res
}
