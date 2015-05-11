// This is just a test run to see how Go works with http services.
// I might do something more interesting than the standard wikia.
// Since I am going to work with MySQL, I think it would be
// appropriate to generate some data and run over it with some
// algorithm and see if we can get something close to the function
// that generated it.
package main

import (
//	"html/template"
//	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
//	"regexp"
	"bytes"
	"encoding/json"
)

type Point struct {
	ID int
	X int
	Y int
}

// generatePoint takes in an id so that when we add it to the database,
// we can have the primary key be the id and since it is incremented,
// inserts should be blazingly fast.
func generatePoint(id int) Point {
	x := int(rand.Int31() / 2)
	y := x * 2
	return Point{ id, x, y }
}

// createJsonPackage takes in an id and creates a json to send to the server
func createJsonPackage(id int) ([]byte, error){
	point := generatePoint(id)
	mar, err := json.Marshal(point)
	if err != nil {
		log.Println(err)
	}
	var p Point
	err = json.Unmarshal(mar, &p)
	if err != nil {
		log.Println(err)
	}
	
	return json.Marshal(p)
}

// postToserver fires up a client and sends the json to the server to be
// unmarshaled and send to the db.
func postToServer(url string, outCount int) error {
	client := &http.Client{}
	i := 0
	var err error
	for i < 1000 {
		pay, err := createJsonPackage(outCount * 100000 + i)
		if err != nil {
			break
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(pay))
		if err != nil {
			break
		}

		resp, err := client.Do(req)
		if err != nil {
			break
		}
		defer resp.Body.Close()

		if err != nil {
			break
		}
		defer resp.Body.Close()
		
		if resp.StatusCode >= 400 {
			break
		}
		i = i + 1
	}
	return err
}

// The program is really just a tester to see if the server works.
func main() {
	url := "http://127.0.0.1:8080"
	for i := 0; i < 10; i++ {
		if err := postToServer(url,i); err != nil {
			log.Println(err)
		}
	}
}
