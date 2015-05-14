package main

import (
	"testing"
	"net/http"
	"log"
	"encoding/json"
	"math/rand"
	"bytes"
	"os"
	//"net/http"
	"net/http/httptest"
	//"database/sql"
	_ "github.com/erikstmartin/go-testdb"
)

var (
	db, er = SetupDatabase()
)

// generatePoint takes in an id so that when we add it to the database,
// we can have the primary key be the id and since it is incremented,
// inserts should be blazingly fast.
func generatePoint(id int) Point {
	x := int(rand.Int31() / 2)
	y := x * 2
	return Point{id, x, y}
}

// createJsonPackage takes in an id and creates a json to send to the server
func createJsonPackage(id int) ([]byte, error) {
	point := generatePoint(id)
	pay, err := json.Marshal(point)
	if err != nil {
		return nil, err
	}
	return pay, nil
}

func RequestGenerator(t *testing.T) *http.Request {
	http.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewReader(pay))
}

func TestHandler(t *testing.T) {	
	pay, err := createJsonPackage(1)
	if err != nil {
		t.Fatal("The Json could not be Marshaled up")
	}

	reqs := [2]*http.Request{
		http.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewReader(pay)),
		http.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewReader(pay)),
	}
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	Handler(w, r, db)
	if w.Code != 200 {
		t.Errorf("%d - %s", w.Code, w.Body.String())
	}
}

func TestMain(m *testing.M) {
	if er != nil {
		log.Fatal("The database could not be instantiated")
	}
	defer db.Close()

	os.Exit(m.Run())
}
