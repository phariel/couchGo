package main

import (
	"encoding/json"
	"fmt"

	"./couch"
)

type testData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fmt.Println("start")
	testJSON, err := json.Marshal(&testData{Name: "test", Age: 11})
	if err == nil {
		res := couch.Insert(testJSON, "person_detect")
		fmt.Println(res)
	}
}
