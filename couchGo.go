package couchGo

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

//Private

func init() {
	srv = &CouchServer{ip: "127.0.0.1", port: 5984, defaultDbName: "master"}
}

type couchRequest struct {
	method    string
	dbName    string
	id        string
	jsonBody  []byte
	timestamp time.Time
}

func (c *couchRequest) doRequest(ch chan ResponseData) {
	if c.dbName != "" {
		urlArray := []string{"http://", srv.ip, ":", strconv.Itoa(srv.port), "/", c.dbName, "/", c.id}
		url := strings.Join(urlArray, "")
		req, _ := http.NewRequest(c.method, url, bytes.NewBuffer(c.jsonBody))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			ch <- ResponseData{Err: false, Status: resp.StatusCode, Body: buf.String()}
		}
	} else {
		ch <- ResponseData{Err: true, Status: 403, Body: "{\"error\":true}"}
	}
}

//Public
var srv *CouchServer

type CouchServer struct {
	ip            string
	port          int
	defaultDbName string
}

type ResponseData struct {
	Err    bool
	Status int
	Body   string
}

type RequestData struct {
	DbName   string
	Id       string
	JsonBody []byte
}

func ConfigServer(cs *CouchServer) {
	srv = cs
}

func Insert(rBatch []RequestData, cb func(res []ResponseData)) {
	var res []ResponseData
	ch := make(chan ResponseData, len(rBatch))
	for _, v := range rBatch {
		func() {
			if v.DbName == "" {
				v.DbName = srv.defaultDbName
			}

			if v.Id == "" {
				v.Id = uuid.NewV4().String()
			}

			cr := couchRequest{method: "PUT", dbName: v.DbName, id: v.Id, jsonBody: v.JsonBody, timestamp: time.Now()}

			go cr.doRequest(ch)
		}()
	}
	batchRequestsNotEnd := true
	for batchRequestsNotEnd {
		if len(ch) == len(rBatch) {
			batchRequestsNotEnd = false
			for i := 0; i < len(rBatch); i++ {
				res = append(res, <-ch)
			}
			cb(res)
		}
	}
}
