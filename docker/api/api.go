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
	"fmt"
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

	pattern := "^April|^May|^Raft|^Blockchain"
	filter := r.URL.Query().Get("filter")

	fmt.Println("Receiving request for name = " + name + " and property = " + prop)

	if filter != "" {
		pattern = filter
	}

	session, err := mgo.Dial("mongodb")
	//session, err := mgo.Dial("54.186.74.114")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("blockchain0").C(name)

	//property := "$messages_count"
	property := "$" + prop

	//match := m{"$match": m{"name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}}}
	//
	//if "block_valid_ratio_percentage" == prop {
	//	fmt.Println(" ... Adding filter!")
	//	match = m{"$match": m{"name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}, "block_valid_ratio_percentage": m{"$gt": 10}}}
	//}

	pipeLine := []m{
		//{"$match": m{"block_valid_ratio_percentage": m{"$gt": 10}}},
		{"$match": m{"block_valid_ratio_percentage": m{"$gt": 10}, "name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}}},
		//{"$match": m{"name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}}},
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
		{"$sort": m{"_id": 1}},
	}

	var result []AggregationResult
	c.Pipe(pipeLine).All(&result)

	if prop == "average_medium_time" || prop == "buffer_channel_time" {
		divValue := float64(1000000)
		for _, element := range result {
			element.MinimumValue = element.MinimumValue / divValue
			element.MaximumValue = element.MaximumValue / divValue
			element.AverageValue = element.AverageValue / divValue
			element.StandardDeviationPValue = element.StandardDeviationPValue / divValue
			element.StandardDeviationSValue = element.StandardDeviationSValue / divValue
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
