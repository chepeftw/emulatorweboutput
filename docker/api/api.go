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
	"math"
	"strconv"
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
	Speed    float64 `bson:"speed"`

	Runs float64 `bson:"runs"`
}

type Highchart struct {
	Name string    `json:"name"`
	Data []int `json:"data"`
}

type Highcharts struct {
	Highchart []Highchart
}

func GetProcessedProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])
	prop := strings.ToLower(params["prop"])

	pattern := "^April|^May|^Raft|^Blockchain"
	filter := r.URL.Query().Get("filter")
	db := r.URL.Query().Get("db")

	fmt.Println("Receiving request for name = " + name + " and property = " + prop)

	if filter != "" {
		pattern = filter
	}

	dbName := "blockchain0"
	if db != "" {
		dbName = db
	}

	session, err := mgo.Dial("mongodb")
	//session, err := mgo.Dial("54.186.74.114")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(dbName).C(name)

	//property := "$messages_count"
	property := "$" + prop

	match := m{"$match": m{"name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}}}

	if "block_valid_ratio_percentage" == prop {
		fmt.Println("block_valid_ratio_percentage filter!")
		match = m{"$match": m{prop: m{"$gt": 10}, "name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}}}
	} else if "query_complete_ms" == prop {
		fmt.Println("query_complete_ms filter!")
		match = m{"$match": m{"name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}, prop: m{"$lt": 50000, "$gt": 0}}}
	} else if "completion_time_ms" == prop {
		fmt.Println("completion_time_ms filter!")
		match = m{"$match": m{"name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}, prop: m{"$lt": 50000, "$gt": 0}}}
	}

	pipeLine := []m{
		//{"$match": m{"block_valid_ratio_percentage": m{"$gt": 10}}},
		match,
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

func GetProcessedGraph(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])
	prop := strings.ToLower(params["prop"])
	timeout := strings.ToLower(params["tmout"])
	speed := strings.ToLower(params["speed"])
	//name := "monitor_query_complete"
	//prop := "query_complete_ms"
	//timeout := 200
	//speed := 2

	pattern := "^Blockchain"
	db := r.URL.Query().Get("db")

	fmt.Println("Receiving request for graph for name = " + name + " and property = " + prop)

	dbName := "blockchain0"
	if db != "" {
		dbName = db
	}

	session, err := mgo.Dial("mongodb")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(dbName).C(name)

	property := "$" + prop

	group := m{"$group":
	m{"_id": "$name",
		"minVal": m{"$min": property},
		"maxVal": m{"$max": property},
		"avgVal": m{"$avg": property},
		"stdPVal": m{"$stdDevPop": property},
		"stdSVal": m{"$stdDevSamp": property},
		"timeout": m{"$avg": "$timeout"},
		"nodes": m{"$avg": "$nodes"},
		"size": m{"$avg": "$size"},
		"speed": m{"$avg": "$speed"},
		"duration": m{"$avg": "$duration"},
		"runs": m{"$sum": 1}},
	}
	sort := m{"$sort": m{"_id": 1}}

	numberOfNodes := [4]int{20, 30, 40, 50}

	finalResult := Highcharts{}
	finalResult.Highchart = append(finalResult.Highchart, Highchart{"Low density", []int{}})
	finalResult.Highchart = append(finalResult.Highchart, Highchart{"Medium density", []int{}})
	finalResult.Highchart = append(finalResult.Highchart, Highchart{"High density", []int{}})

	for _, nodes := range numberOfNodes {
		fmt.Println("Querying for numberOfNodes = " + strconv.Itoa(nodes))
		match := m{"$match": m{"name": m{"$regex": bson.RegEx{Pattern: pattern, Options: "si"}}, prop: m{"$lt": 50000, "$gt": 0}, "timeout": timeout, "speed": speed, "nodes": nodes}}
		pipeLine := []m{match, group, sort}

		var result []AggregationResult
		c.Pipe(pipeLine).All(&result)

		sum := float64(0)
		min := float64(1000)
		max := float64(0)
		for _, res := range result {
			min = math.Min(min, res.Size)
			max = math.Max(max, res.Size)
			sum += res.Size
		}
		med := sum - min - max

		minVal := 0
		medVal := 0
		maxVal := 0
		for _, res := range result {
			fmt.Println("Checking average value for " + strconv.Itoa(int(res.Size)))
			formatValue := int(res.AverageValue)
			switch res.Size {
			case min:
				minVal = formatValue
				break
			case med:
				medVal = formatValue
				break
			case max:
				maxVal = formatValue
				break
			}
		}

		finalResult.Highchart[0].Data = append(finalResult.Highchart[0].Data, minVal)

		finalResult.Highchart[1].Data = append(finalResult.Highchart[1].Data, medVal)

		finalResult.Highchart[2].Data = append(finalResult.Highchart[2].Data, maxVal)
	}

	json.NewEncoder(w).Encode(finalResult)
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/property/{name}/{prop}", GetProcessedProperty).Methods("GET")
	//router.HandleFunc("/graph/{name}/{prop}", GetProcessedGraph).Methods("GET")
	router.HandleFunc("/graph/{name}/{prop}/{tmout}/{speed}", GetProcessedGraph).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
