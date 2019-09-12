package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/narqo/go-dogstatsd-parser"
)

var widestNameSeen int
var widestTagsSeen int

// Keep track of how our data should be displayed
// TODO: Check terminal width and space appropriately
func init() {
	widestNameSeen = 0
	widestTagsSeen = 0
}

// Adapted from https://stackoverflow.com/q/43947363
func stdoutIsTerminal() bool {
	fi, err := os.Stdout.Stat()

	if err != nil {
		return false
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	return true
}

func printTags(tags map[string]string) string {
	keys := make([]string, 0, len(tags))

	// Copy over names and sort them
	for key, value := range tags {
		keys = append(keys, key+":"+value)
	}

	sort.Strings(keys)

	return strings.Join(keys, ",")
}

func printMetricForTerminal(data *dogstatsd.Metric) {
	// Update widths
	if len(data.Name) > widestNameSeen {
		widestNameSeen = len(data.Name)
	}

	// TODO: Stringify tags uniformly
	tags := printTags(data.Tags)

	if len(tags) > widestTagsSeen {
		widestTagsSeen = len(tags)
	}

	fmtString := fmt.Sprintf("%%-%ds\t%%s\t%%.2f\t%%-%ds\t", widestNameSeen, widestTagsSeen)

	switch data.Value.(type) {
	case float32, float64:
		fmtString += "%0.2f\n"
	case int, int32, int64:
		fmtString += "%d\n"
	default:
		fmtString += "%v\n"
	}

	// TODO: Drop in a timestamp...
	fmt.Printf(fmtString, data.Name, data.Type, data.Rate, tags, data.Value)

}

func printMetricForCharDevice(data *dogstatsd.Metric) {
	tags := printTags(data.Tags)

	switch data.Value.(type) {
	case float32, float64:
		fmt.Printf("%s\t%s\t%.2f\t%s\t%0.2f\n", data.Name, data.Type, data.Rate, tags, data.Value)
	case int, int32, int64:
		fmt.Printf("%s\t%s\t%.2f\t%s\t%d\n", data.Name, data.Type, data.Rate, tags, data.Value)
	default:
		fmt.Printf("%s\t%s\t%.2f\t%s\t%v\n", data.Name, data.Type, data.Rate, tags, data.Value)
	}
}

func main() {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 8125,
		IP:   net.ParseIP("127.0.0.1"),
	})

	if err != nil {
		panic(err)
	}

	printer := printMetricForTerminal

	if stdoutIsTerminal() == false {
		printer = printMetricForCharDevice
	}

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
				printer(metrics[i])
			}
		}
	}
}
