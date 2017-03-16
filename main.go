package main

import (
	"fmt"
	"net"
	"strings"
	"sort"

	"github.com/narqo/go-dogstatsd-parser"
)

var widestNameSeen int
var widestTagsSeen int

func printTags(tags map[string]string) string {
	keys := make([]string, 0, len(tags))

	// Copy over names and sort them
	for key,value := range tags {
		keys = append(keys, key + ":" + value)
	}

	sort.Strings(keys)

	return strings.Join(keys, ",")
}

func displayMetric(data *dogstatsd.Metric) {
	// Update widths
	if len(data.Name) > widestNameSeen {
		widestNameSeen = len(data.Name)
	}

	// TODO: Stringify tags uniformly
	tags := printTags(data.Tags)

	if len(tags) > widestTagsSeen {
		widestTagsSeen = len(tags)
	}

	fmtString := fmt.Sprintf("%%-%ds  %%-%ds ", widestNameSeen, widestTagsSeen)

	switch data.Value.(type) {
	case float32, float64:
		fmtString += "%0.2f\n"
	case int, int32, int64:
		fmtString += "%d\n"
	default:
		fmtString += "%v\n"
	}

	// TODO: Drop in a timestamp...
	fmt.Printf(fmtString, data.Name, tags, data.Value)

}

func main() {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 8125,
		IP:   net.ParseIP("127.0.0.1"),
	})

	if err != nil {
		panic(err)
	}

	// Keep track of how our data should be displayed
	// TODO: Check terminal width and space appropriately
	widestNameSeen = 0
	widestTagsSeen = 0

	data := make([]byte, 65535)
	for {
		length, _, err := ln.ReadFromUDP(data)

		//fmt.Println(length, addr, err, string(data[0:length]))

		lines := strings.Fields(string(data[0:length]))
		metrics := make([]*dogstatsd.Metric, len(lines))

		//fmt.Printf("data: %+v\n", lines)

		for i, line := range lines {
			metrics[i], err = dogstatsd.Parse(line)
			if err != nil {
				fmt.Printf("ERR: %s -> %+v\n", line, err)
			} else {
				displayMetric(metrics[i])
			}
		}
	}
}
