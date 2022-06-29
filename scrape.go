package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Option struct {
	code string
	name string
}

type Attribute struct {
	name string
	value string
}

type USAirNetMapStationData struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type USAirNetMapStateData struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Stations map[string]*USAirNetMapStationData `json:"stations"`
}

type USAirNetMapData struct {
	States map[string]*USAirNetMapStateData `json:"states"`
}

var data USAirNetMapData

func getSelectOptions(node *html.Node) []Option {
	options := make([]Option, 0, 1000)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "option" {
			var val string
			for _, a := range c.Attr {
				if a.Key == "value" {
					val = a.Val
				}
			}
			if val != "" {
				options = append(options, Option{val, c.FirstChild.Data})
			}
		}
	}
	return options
}


func getElement(node *html.Node, tag string, attrib Attribute) (*html.Node, bool) {
	if node.Data == tag {
		for _, a := range node.Attr {
			if a.Key == attrib.name && a.Val == attrib.value {
				return node, false
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		subnode, fail := getElement(c, tag, attrib)
		if !fail {
			return subnode, false
		}
	}
	return nil, true
}

func getHTML(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(bodyBytes)
}

func getNodeFromUrl(url string) *html.Node {
	body := getHTML(url)
	//fmt.Println(body)
	root, err := html.Parse(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	return root
}

func scrapeStation(state string, station string) {
	url := "http://www.usairnet.com/cgi-bin/launch/code.cgi?Submit=Go&state=" + state + "&sta=" + station
	//fmt.Println(url)
	//body := getHTML(url)
	//fmt.Println(body)
	root := getNodeFromUrl(url)

	// get the heading that contains the station title
	headingSpan, fail := getElement(root, "span", Attribute{"class", "bolder"})
	if fail {
		log.Fatal("scrapeStation(" + state + ", " + station + ") failed to get span class bolder")
	}
	headingStrong := headingSpan.FirstChild
	headingText := headingStrong.FirstChild.Data
	//fmt.Printf("headingText: %s\n", headingText)

	headingPrefix := "Aviation Weather Forecast at "
	if !strings.Contains(headingText, headingPrefix) {
		log.Fatal("Heaing text (" + headingText + ") doesn't contain expected prefix (" + headingPrefix + ")")
	}

	stationText := strings.Replace(headingText, headingPrefix, "", -1)
	//fmt.Printf("stationText: %s\n", stationText)

	//stationText = "Franklin, Somewhere, Pennsylvania"
	// intentionally split by comma and space, then rejoin for the title
	stationTextParts := strings.Split(stationText, ", ")
	stationName := strings.Join(stationTextParts[0:len(stationTextParts)-1], ", ")
	stateName := stationTextParts[len(stationTextParts)-1]
	//fmt.Printf("stationName: %s\n", stationName)
	//fmt.Printf("stateName: %s\n", stateName)

	// get the details line, and extract out station code (for confirmation) and coords
	detailSpan, fail := getElement(root, "span", Attribute{"class", "norm2"})
	if fail { log.Fatal(fail) }

	text := detailSpan.FirstChild.Data
	spaces := regexp.MustCompile(` +`)

	tokens := spaces.Split(text, -1)

	if tokens[0] != "Station:" { log.Fatal("Error finding 'Station:' in " + text) }
	if tokens[1] != station { log.Fatal("Scraped station code " + tokens[1] + ", expected " + station) }

	// all USA should be North/West, so just assert them
	if tokens[3] != "North:" { log.Fatal("Error finding 'North:' name in " + text) }
	north, err := strconv.ParseFloat(tokens[4], 32)
	if err != nil { log.Fatal(err) }

	if tokens[6] != "West:" {
		log.Fatal("Error finding 'West:' name in " + text)
	}
	west, err := strconv.ParseFloat(tokens[7], 32)
	if err != nil { log.Fatal(err) }

	// write scraped information to <data>
	//- ensure states map exists
	if data.States == nil { data.States = map[string]*USAirNetMapStateData{} }
	//- ensure stateData object exists
	stateData := data.States[state]
	if stateData == nil {
		data.States[state] = &USAirNetMapStateData{}
		stateData = data.States[state]
	}
	//- assign state code if not yet assigned
	if stateData.Code == "" { stateData.Code = state }
	//- assign state name if first time, or confirm
	if stateData.Name == "" {
		stateData.Name = stateName
		fmt.Println(stateName)
	} else {
		if stateData.Name != stateName {
			log.Fatal("scraped state name " + stateName + " not equal to existing state name " + stateData.Name)
		}
	}
	//- ensure stations is not nil
	if stateData.Stations == nil { stateData.Stations = map[string]*USAirNetMapStationData{} }
	//- ensure stationData object exists
	stationData := stateData.Stations[station]
	if stationData == nil {
		stateData.Stations[station] = &USAirNetMapStationData{}
		stationData = stateData.Stations[station]
	}
	//- assign station code if not yet assigned
	if stationData.Code == "" { stationData.Code = station }
	//- assign station name if first time, or confirm
	if stationData.Name == "" {
		stationData.Name = stationName
		fmt.Println(stationName)
	} else {
		if stationData.Name != stationName {
			log.Fatal("scraped station name " + stationName + " not equal to existing station name " + stationData.Name)
		}
	}
	//- assign values for North/West without checking
	stationData.Latitude = north
	stationData.Longitude = west

	fmt.Printf("%s Station %s was successfully scraped\n", state, station)
	fmt.Printf("- Name: %s\n", stationName)
	fmt.Printf("- Latitude: %f North\n", north)
	fmt.Printf("- Longitude: %f West\n", west)
	fmt.Printf("Saving...")
	fileData, _ := json.MarshalIndent(data, "","\t")
	_ = ioutil.WriteFile("data.json", fileData, 0644)
	fmt.Printf(" done!\n")
}


func scrapeState(state string) {
	url := "http://www.usairnet.com/cgi-bin/launch/code.cgi?Submit=Go&state=" + state
	root := getNodeFromUrl(url)
	sel, fail := getElement(root, "select", Attribute{"name", "sta"})
	if fail {
		log.Fatal("scrapeState(" + state + ") failed to get 'select' element named 'sta'")
	}
	stations := getSelectOptions(sel)
	//fmt.Println(len(stations))
	for i, station := range stations {
		fmt.Printf("station %d: %s - %s\n", i, station.code, station.name)
		scrapeStation(state, station.code)
		fmt.Printf("\nSleeping for 2secs to give website a break... ")
		time.Sleep(2 * time.Second)
		fmt.Printf("ok let's go again!\n\n")
	}
	fmt.Printf("Done with %s\n", state)
}

func scrapeUSA() {
	url := "http://www.usairnet.com/cgi-bin/launch/code.cgi"
	root := getNodeFromUrl(url)
	sel, fail := getElement(root, "select", Attribute{"name", "state"})
	if fail {
		log.Fatal("scrapeUSA() failed to get 'select' element named 'state'")
	}
	states := getSelectOptions(sel)
	//fmt.Println(len(states))
	for i, state := range states {
		fmt.Printf("state %d: %s - %s\n", i, state.code, state.name)
		scrapeState(state.code)
		fmt.Printf("\nSleeping for 2secs to give website a break... ")
		time.Sleep(2 * time.Second)
		fmt.Printf("ok let's go again!\n\n")
	}
}

func main() {
	args := os.Args
	//fmt.Println(reflect.TypeOf(args))
	//fmt.Println(len(args))

	file, _ := ioutil.ReadFile("data.json")
	json.Unmarshal([]byte(file), &data)

	/*
	states := data.States
	for state, stateData := range states {
		fmt.Printf("states[%s]: %s\n", state, stateData.Name)
		for station, stationData := range stateData.Stations {
			fmt.Printf("states[%s].stations[%s]: %s\n", state, station, stationData.Name)
		}
	}
	*/

	if len(args) == 1 {
		scrapeUSA()
	} else if len(args) == 2 {
		scrapeState(args[1])
	} else if len(args) == 3 {
		scrapeStation(args[1], args[2])
	} else {
		fmt.Println("Usage: go run scrape.go [state] [station]")
	}
}
