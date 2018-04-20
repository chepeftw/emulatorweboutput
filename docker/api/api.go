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
)

type m bson.M // just for brevity, bson.M type is map[string]interface{}

type AggregationResult struct {
	Id                      string  `bson:"_id"`
	MinimumValue            float64 `bson:"minVal"`
	MaximumValue            float64 `bson:"maxVal"`
	AverageValue            float64 `bson:"avgVal"`
	StandardDeviationPValue float64 `bson:"stdPVal"`
	StandardDeviationSValue float64 `bson:"stdSVal"`

	Timeout  float64 `bson:"timeout"`
	Nodes    float64 `bson:"nodes"`
	Size     float64 `bson:"size"`
	Duration float64 `bson:"duration"`

	Runs float64 `bson:"runs"`
}

func GetProcessedProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])
	prop := strings.ToLower(params["prop"])

	session, err := mgo.Dial("mongodb")
	//session, err := mgo.Dial("54.186.74.114")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("blockchain").C(name)

	//property := "$messages_count"
	property := "$" + prop

	pipeLine := []m{
		{"$match": m{"name": m{ "$regex": bson.RegEx{Pattern: `^April|^Raft`, Options: "si"} }}},
		{"$group":
		m{"_id": "$name",
			"minVal": m{"$min": property},
			"maxVal": m{"$max": property},
			"avgVal": m{"$avg": property},
			"stdPVal": m{"$stdDevPop": property},
			"stdSVal": m{"$stdDevSamp": property},
			"timeout": m{"$avg": "$timeout"},
			"nodes": m{"$avg": "$nodes"},
			"size": m{"$avg": "$size"},
			"duration": m{"$avg": "$duration"},
			"runs": m{"$sum": 1}}},
	}

	var result []AggregationResult
	c.Pipe(pipeLine).All(&result)

	if prop == "average_medium_time" || prop == "buffer_channel_time" {
		for _, element := range result {
			element.MinimumValue = element.MinimumValue/1000000
			element.MaximumValue = element.MaximumValue/1000000
			element.AverageValue = element.AverageValue/1000000
			element.StandardDeviationPValue = element.StandardDeviationPValue/1000000
			element.StandardDeviationSValue = element.StandardDeviationSValue/1000000
		}
	}

	json.NewEncoder(w).Encode(result)
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/property/{name}/{prop}", GetProcessedProperty).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
