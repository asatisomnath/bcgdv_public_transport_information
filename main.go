// main.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Delay struct {
	LineName string `json:"line_name"`
	Delay    string `json:"delay"`
}

type Time struct {
	LineId string `json:"line_id"`
	StopId string `json:"stop_id"`
	Time   string `json:"time"`
}

type Line struct {
	LineId   string `json:"line_id"`
	LineName string `json:"line_name"`
}

type Stop struct {
	StopId string `json:"stop_id"`
	X      string `json:"x"`
	Y      string `json:"y"`
}

var Delays []Delay

var Times []Time

var Lines []Line

var Stops []Stop


func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllDelays(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllDelays")
	json.NewEncoder(w).Encode(Delays)
}

func returnSingleDelay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, delay := range Delays {
		if delay.LineName == key {
			json.NewEncoder(w).Encode(delay)
		}
	}
}

func createNewDelay(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Delay struct
	// append this to our Delays array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var delay Delay
	json.Unmarshal(reqBody, &delay)
	// update our global Articles array to include
	// our new Article
	Delays = append(Delays, delay)

	json.NewEncoder(w).Encode(delay)
}

func deleteDelay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, delay := range Delays {
		if delay.LineName == id {
			Delays = append(Delays[:index], Delays[index+1:]...)
		}
	}

}

func returnArrivingVehicleByStop(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["stopId"]

	json.NewEncoder(w).Encode(returnArrivingVehicle(key, time.Now()))

}

func returnArrivingVehicleByLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["time"]
	x := vars["x"]
	y := vars["y"]

	t, _ := time.Parse("15:04:05", key)

	for _, stop := range Stops {
		if stop.X == x && stop.Y == y {
			json.NewEncoder(w).Encode(returnArrivingVehicle(stop.StopId, t))
			return
		}
	}

	var stopId = ""

	var nearest = math.MaxFloat64

	for _, stop:= range Stops{

		x1 , _ :=strconv.Atoi(x)
		x2 , _ :=strconv.Atoi(stop.X)
		y1 , _ :=strconv.Atoi(y)
		y2 , _ :=strconv.Atoi(stop.Y)

		var val = math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2))

		if val<nearest {
			nearest = val
			stopId = stop.StopId
		}
	}


	json.NewEncoder(w).Encode(returnArrivingVehicle(stopId, t))

}

func returnArrivingVehicle(stopId string, cur time.Time) Time {

	var lineName = ""
	var firstLine = ""
	latest := cur

	var isExist = false
	var isFirstLine = true
	for _, t := range Times {
		if t.StopId == stopId {
			t1, _ := time.Parse("15:04:05", t.Time)

			for _, line := range Lines {
				if t.LineId == line.LineId {
					if isFirstLine {
						isFirstLine = false
						firstLine = line.LineName
						latest = t1
					}
					lineName = line.LineName
					for _, delay := range Delays {
						if delay.LineName == line.LineName {
							count, _ := strconv.Atoi(delay.Delay)

							t1 = t1.Add(time.Duration(count) * time.Minute)
						}
					}

				}
			}

			if !isExist && t1.After(cur) {
				isExist = true
				latest = t1
				firstLine = lineName
			} else if isExist && t1.Before(latest) && t1.After(cur) {
				latest = t1
				firstLine = lineName
			}
		}
	}

	return Time{ LineId: firstLine, StopId: stopId, Time: latest.Format("15:04:05")}
}


func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/delays", returnAllDelays)
	myRouter.HandleFunc("/delay", createNewDelay).Methods("POST")
	myRouter.HandleFunc("/delay/{id}", deleteDelay).Methods("DELETE")
	myRouter.HandleFunc("/delay/{id}", returnSingleDelay)

	myRouter.HandleFunc("/arriving/{stopId}", returnArrivingVehicleByStop)
	myRouter.HandleFunc("/arriving/{time}/{x}/{y}", returnArrivingVehicleByLocation)


	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {

	csvReaderDelays()
	csvReaderLines()
	csvReaderTimes()
	csvReaderStops()

	handleRequests()
}

func csvReaderDelays() {
	// Open the file
	recordFile, err := os.Open("delays.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	_, err = reader.Read()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	Delays = []Delay{}

	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			fmt.Println("An error encountered ::", err)
			return
		}
		Delays = append(Delays, Delay{LineName: record[0], Delay: record[1]})
	}

	err = recordFile.Close()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}
}

func csvReaderLines() {
	// Open the file
	recordFile, err := os.Open("lines.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	_, err = reader.Read()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	Lines = []Line{}

	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			fmt.Println("An error encountered ::", err)
			return
		}
		Lines = append(Lines, Line{LineId: record[0], LineName: record[1]})
	}

	err = recordFile.Close()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}
}

func csvReaderTimes() {
	// Open the file
	recordFile, err := os.Open("times.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	_, err = reader.Read()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	Times = []Time{}

	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			fmt.Println("An error encountered ::", err)
			return
		}
		Times = append(Times, Time{LineId: record[0], StopId: record[1], Time: record[2]})
	}

	err = recordFile.Close()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}
}



func csvReaderStops() {
	// Open the file
	recordFile, err := os.Open("stops.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	_, err = reader.Read()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	Stops = []Stop{}

	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			fmt.Println("An error encountered ::", err)
			return
		}
		Stops = append(Stops, Stop{StopId: record[0], X: record[1], Y: record[2]})
	}

	err = recordFile.Close()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}
}
