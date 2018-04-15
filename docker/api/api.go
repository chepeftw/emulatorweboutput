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

type m bson.M // just for brevity, bson.M type is map[string]interface{}

type EmulationResult struct {
	Url string
	Name string
	Created int64

	Result1 string
}

type AggregationResult struct {
	Id        string `bson:"_id"`
	MinimumValue float64 `bson:"minVal"`
	MaximumValue float64 `bson:"maxVal"`
	AverageValue float64 `bson:"avgVal"`
	StandardDeviationPValue float64 `bson:"stdPVal"`
	StandardDeviationSValue float64 `bson:"stdSVal"`

	Timeout float64 `bson:"timeout"`
	Nodes float64 `bson:"nodes"`
	Size float64 `bson:"size"`
	Duration float64 `bson:"duration"`

	Runs float64 `bson:"runs"`
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
	err = c.Find(bson.M{"url": "http://emulator.chepeftw.com/"}).Sort("-created").Limit(1).One(&results)

	json.NewEncoder(w).Encode(results)
}

func GetProcessedProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])
	prop := strings.ToLower(params["prop"])

	session, err := mgo.Dial("mongodb")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("blockchain").C(name)

	//property := "$messages_count"
	property := "$"+prop

	pipeLine := []m{
		m{"$match": m{"name": "JulyTest_20_0"}},
		m{"$group":
			m{ "_id": "$name",
			"minVal": m{ "$min": property },
			"maxVal": m{ "$max": property },
			"avgVal": m{ "$avg": property },
			"stdPVal": m{ "$stdDevPop": property },
			"stdSVal": m{ "$stdDevSamp": property },
			"timeout": m{ "$avg": "$timeout" },
			"nodes": m{ "$avg": "$nodes" },
			"size": m{ "$avg": "$size" },
			"duration": m{ "$avg": "$duration" },
			"runs": m{ "$sum": 1 } }},
	}
	var result AggregationResult
	//result := AggregationResult{}
	c.Pipe(pipeLine).One(&result)

	json.NewEncoder(w).Encode(result)
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/property/{name}/{prop}", GetProcessedProperty).Methods("GET")
	router.HandleFunc("/test/set/{name}", SetTestData).Methods("GET")
	router.HandleFunc("/test/get/{name}", GetTestData).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}