package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"github.com/rs/cors"
	"strings"
	"time"
)

type EmulationResult struct {
	Url string
	Created int64

	Result1 string
}

// Display a single data
func GetICO(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	session, err := mgo.Dial("mongodb")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	name := params["name"]
	name = strings.ToLower(name)

	//c := session.DB("emulation").C("version1")

	results := EmulationResult{
		Url: "https://semantic-ui.com/examples/grid.html",
		Created: time.Now().Unix(),
		Result1: "123456789",
	}
	//err = c.Find(bson.M{"name": name}).Sort("-created").Limit(1).All(&results)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(results)
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/emulation/{name}", GetICO).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}