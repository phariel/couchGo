package couchGo

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

//RequestData --
type RequestData struct {
	DbName   string
	ID       string
	JSONBody []byte
}

// ResponseData --
type ResponseData struct {
	ID          string
	Err         bool
	Status      int
	Body        string
	RequestTime time.Time
}

type couchRequest struct {
	method   string
	dbName   string
	id       string
	jsonBody []byte
	time     time.Time
}

func (c *couchRequest) doRequest(ch chan ResponseData) {
	if c.dbName != "" && c.id != "" {
		urlArray := []string{"http://", srv.ip, ":", strconv.Itoa(srv.port), "/", c.dbName, "/", c.id}
		url := strings.Join(urlArray, "")
		req, _ := http.NewRequest(c.method, url, bytes.NewBuffer(c.jsonBody))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			ch <- ResponseData{ID: c.id, Err: false, Status: resp.StatusCode, Body: buf.String(), RequestTime: c.time}
		}
	} else {
		ch <- ResponseData{ID: c.id, Err: true, Status: 403, Body: "{\"error\":true}", RequestTime: c.time}
	}
}

func batchRequests(method string, rBatch []RequestData, cb func(res []ResponseData)) {
	var res []ResponseData
	ch := make(chan ResponseData, len(rBatch))
	for _, v := range rBatch {
		func() {
			if v.DbName == "" {
				v.DbName = srv.defaultDbName
			}

			if method == "PUT" && v.ID == "" {
				v.ID = uuid.NewV4().String()
			}

			cr := couchRequest{method: method, dbName: v.DbName, id: v.ID, jsonBody: v.JSONBody, time: time.Now()}

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
