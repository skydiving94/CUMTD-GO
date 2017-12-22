package main

import (
	"fmt"
	"os"
	"net/http"
//	"strconv"
	"encoding/json"
	"io/ioutil"
//	"strings"
//	"flag"
	"time"
)

type Route_t struct {
	RouteColor string `json:"route_color"`
	RouteId string `json:"route_id"`
	RouteLongName string `json:"route_long_name"`
	RouteShortName string `json:"route_short_name"`
	RouteTextColor string `json:"route_text_color"`
}

type Trip_t struct {
	TripId string `json:"trip_id"`
	TripHeadsign string `json:"trip_headsign"`
	RouteId string `json:"route_id"`
	BlockId string `json:"block_id"`
	Direction string `json:"direction"`
	ServiceId string `json:"service_id"`
	ShapeId string `json:"shape_id"`
}

type StopTime_t struct {
	ArrivalTime string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
	StopSequence int32 `json:"stop_sequence"`
	StopId string `json:"stop_id"`
	Trip Trip_t `json:"trip"`
}

type StopPoint_t struct {
	Code string `json:"code"`
	StopID string `json:"stop_id"`
	StopLat float64 `json:"stop_lat"`
	StopLon float64 `json:"stop_lon"`
	StopName string `json:"stop_name"`
}

type Stop_t struct {
	StopId string `json:"stop_id"`
	StopName string `json:"stop_name"`
	Code string `json:"code"`
	PercentMatch int64 `json:"percent_match"`
	StopPoints []StopPoint_t `json:"stop_points"`
}

type Departure_t struct {
	StopId string `json:"stop_id"`
	Headsign string `json:"headsign"`
	Route Route_t `json:"route"`
	Trip Trip_t `json:"trip"`
	VehicleId string `json:"vehicle_id"`
	Origin map[string]string `json:"origin"`
	Destination map[string]string `json:"destination"`
	IsMonitored bool `json:"is_monitored"`
	IsScheduled bool `json:"is_scheduled"`
	IsIstop bool `json:"is_istop"`
	Scheduled string `json:"scheduled"`
	Expected string `json:"expected"`
	ExpectedMins int64 `json:"expected_mins"`
	Location map[string]float64 `json:"location"`
}

type Status_t struct {
	Code int32 `json:"code"`
	Msg string `json:"msg"`
}

type Rqst_t struct {
	Method string `json:"method"`
	Params map[string]string `json:"params"`
}

type UnmarshalStopTimesStruct_t struct {
	Time string `json:"time"`
	ChangeSetId string `json:"changeset_id"`
	NewChangeSet bool `json:"new_changeset"`
	Status Status_t `json:"status"`
	Rqst Rqst_t `json:"rqst"`
	StopTimes []StopTime_t `json:"stop_times"`
}

type UnmarshalDeparturesStruct_t struct {
	Time string `json:"time"`
	NewChangeSet bool `json:"new_changeset"`
	Status Status_t `json:"status"`
	Rqst Rqst_t `json:"rqst"`
	Departures []Departure_t `json:"departures"`
}

type UnmarshalStopsStruct_t struct {
	Time string `json:"time"`
	ChangeSetId string `json:"changeset_id"`
	NewChangeSet bool `json:"new_changeset"`
	Status Status_t `json:"status"`
	Rqst Rqst_t `json:"rqst"`
	Stops []Stop_t `json:"stops"`
}

var apis = map[string]string {
	"getStopsBySearch" : "https://developer.cumtd.com/api/v2.2/json/getstopsbysearch?key=eee387f2e3694365a638ee562a9caebd&query=",
	"getStopTimeByStop" : "https://developer.cumtd.com/api/v2.2/json/getstoptimesbystop?key=eee387f2e3694365a638ee562a9caebd&stop_id=", // Simply append a stop id.
	"getDeparturesByStop" : "https://developer.cumtd.com/api/v2.2/json/getdeparturesbystop?key=eee387f2e3694365a638ee562a9caebd&stop_id=",
}

func getStopsBySearch(Query string) (Stops []Stop_t) {
	api := apis["getStopsBySearch"]
	url := api + Query
	
	res, err := http.Get(url)
	if(err != nil) {
		return Stops
	}

	body, err := ioutil.ReadAll(res.Body)
	if(err != nil) {
		return Stops
	}

	var unmarshalStruct UnmarshalStopsStruct_t

	unmarshalErr := json.Unmarshal(body, &unmarshalStruct)
	if unmarshalErr != nil {
		fmt.Println("error:", unmarshalErr)
	}

	return unmarshalStruct.Stops
}

func getStopTime(StopId string, CurrentTime time.Time) (StopTimes []StopTime_t) {
	// hour, min, _ := CurrentTime.Clock()

	api := apis["getStopTimeByStop"]
	url := api + StopId

	res, err := http.Get(url)

	if(err != nil) {
		return StopTimes
	}

	body, err := ioutil.ReadAll(res.Body)
	if(err != nil) {
		return StopTimes
	}

	var unmarshalStruct UnmarshalStopTimesStruct_t

	unmarshalErr := json.Unmarshal(body, &unmarshalStruct)
	if unmarshalErr != nil {
		fmt.Println("error:", unmarshalErr)
	}

	minTime := TimeToSecond(CurrentTime)
	maxTime := minTime + 3600
	
	for i := 0; i < len(unmarshalStruct.StopTimes); i++ {
		arrivalTime := unmarshalStruct.StopTimes[i].ArrivalTime
		arrivalTimeInSecond := StringToSecond(arrivalTime)
		if arrivalTimeInSecond >= minTime && arrivalTimeInSecond <= maxTime {
			StopTimes = append(StopTimes, unmarshalStruct.StopTimes[i])
		}
	}
	
	return StopTimes
}

func getDeparturesByStop(StopId string)(Departures []Departure_t) {
	api := apis["getDeparturesByStop"]

	url := api + StopId

	res, err := http.Get(url)

	if(err != nil) {
		return Departures
	}

	body, err := ioutil.ReadAll(res.Body)
	if(err != nil) {
		return Departures
	}

	var unmarshalStruct UnmarshalDeparturesStruct_t

	unmarshalErr := json.Unmarshal(body, &unmarshalStruct)
	if unmarshalErr != nil {
		fmt.Println("error:", unmarshalErr)
	}

	return unmarshalStruct.Departures
}

func printStops(stops []Stop_t) {
	if len(stops) == 0 {
		fmt.Println("No stops found.")
	} else {
		for i := 0; i < len(stops); i++ {
			fmt.Println("StopId: ", stops[i].StopId)
			fmt.Println("StopName: ", stops[i].StopName)
			fmt.Println("")
		}
	}
}

func printDepartures(departures []Departure_t) {
	if len(departures) == 0 {
		fmt.Println("No departures found.")
	} else {	
		for i := 0; i < len(departures); i++ {
			fmt.Println(departures[i].Headsign)
			fmt.Println("Scheduled: ", departures[i].Scheduled)
			fmt.Println("Expected:  ", departures[i].Expected)
			fmt.Println("")
		}
	}
}

func main() {
//	var currentTime = time.Now()
//	stopTimes := getStopTime("walmart", currentTime)
	//	printTrips(stopTimes)
	args := os.Args

	if len(args) != 3 {
		fmt.Println("Please enter the name of the stop, and Y/N for search or not.")
		return
	}

	if args[2] == "Y" {
		stops := getStopsBySearch(args[1])
		printStops(stops)
		if len(stops) != 0 {
			var stopId string
			fmt.Println("Please enter the stop id for the stop you are looking for.")
			fmt.Scanln(&stopId)
			printDepartures(getDeparturesByStop(stopId))
		}
	} else {
		printDepartures(getDeparturesByStop(args[1]))
	}

	
}
