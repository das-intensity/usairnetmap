package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
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
	//body := getHTML(url)
	//fmt.Println(body)
	root := getNodeFromUrl(url)

	td, fail := getElement(root, "td", Attribute{"class", "norm2"})
	if fail {
		log.Fatal(fail)
	}

	span, fail := getElement(td, "span", Attribute{"class", "norm2"})
	if fail {
		log.Fatal(fail)
	}

	text := span.FirstChild.Data
	spaces := regexp.MustCompile(` +`)
	//fmt.Printf("span: %s\n", text)

	tokens := spaces.Split(text, -1)
	/*
	for _, token := range tokens {
		fmt.Println(token)
	}
	*/

	if tokens[0] != "Station:" {
		log.Fatal("Error finding 'Station:' in " + text)
	}
	if tokens[1] != station {
		log.Fatal("Error finding station name in " + text)
	}

	// all USA should be North/West, so just assert them
	if tokens[3] != "North:" {
		log.Fatal("Error finding 'North:' name in " + text)
	}
	north, err := strconv.ParseFloat(tokens[4], 32)
	if err != nil {
		log.Fatal(err)
	}

	if tokens[6] != "West:" {
		log.Fatal("Error finding 'West:' name in " + text)
	}
	west, err := strconv.ParseFloat(tokens[7], 32)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Station %s successfully found at North: %f, West: %f\n", station, north, west)
}


func scrapeState(state string) {
	url := "http://www.usairnet.com/cgi-bin/launch/code.cgi?Submit=Go&state=" + state
	root := getNodeFromUrl(url)
	sel, fail := getElement(root, "select", Attribute{"name", "sta"})
	if fail {
		log.Fatal(fail)
	}
	stations := getSelectOptions(sel)
	fmt.Println(len(stations))
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
		log.Fatal(fail)
	}
	states := getSelectOptions(sel)
	fmt.Println(len(states))
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
	fmt.Println(reflect.TypeOf(args))
	fmt.Println(len(args))
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
