package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/rs/cors"
	"strings"
	"time"
)

type EmulationResult struct {
	Url string
	Name string
	Created int64

	Result1 string
}

type RandomResponse struct {
	Status string
}

func SetTestData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])

	session, err := mgo.Dial("mongodb")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("emulation").C(name)

	c.Insert(&EmulationResult{
		"http://emulator.chepeftw.com/",
		"Random",
		time.Now().Unix(),
		"yes!" })

	result := RandomResponse{ "OK" }
	json.NewEncoder(w).Encode(result)
}

func GetTestData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])

	session, err := mgo.Dial("mongodb")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("emulation").C(name)

	results := EmulationResult{}
	err = c.Find(bson.M{"url": "http://emulator.chepeftw.com/"}).Sort("-created").Limit(1).All(&results)

	json.NewEncoder(w).Encode(results)
}

// Display a single data
func GetICO(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	name := strings.ToLower(params["name"])

	results := EmulationResult{
		Url: "https://semantic-ui.com/examples/grid.html",
		Name: name,
		Created: time.Now().Unix(),
		Result1: "123456789",
	}

	json.NewEncoder(w).Encode(results)
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/emulation/{name}", GetICO).Methods("GET")
	router.HandleFunc("/test/set/{name}", SetTestData).Methods("GET")
	router.HandleFunc("/test/get/{name}", GetTestData).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}