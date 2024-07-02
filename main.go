package main

import (
	"fmt"
	"net"
	duration_circular_buffer "ping-checker/duration_circular_buffer"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/guptarohit/asciigraph"
)

var latencyRing = duration_circular_buffer.New(100)
var danger = color.New(color.FgRed)
var warning = color.New(color.FgHiYellow)
var success = color.New(color.FgWhite)
var url = "1.1.1.1:80"

func main() {

	writer := uilive.New()
	writer.Start()

	defer writer.Stop()

	for {
		conStatus, latency := checkInternet()
		printColorText(conStatus, latency)
		printLatencyGraph(latency)
		time.Sleep(1 * time.Second)
	}
}

func checkInternet() (bool, time.Duration) {
	start := time.Now()
	con, err := net.DialTimeout("tcp", url, time.Second*1)
	end := time.Now()
	latency := end.Sub(start)
	if err != nil {
		danger.Print("Not connected ")
		danger.Print("\n")
		return false, -1 * time.Millisecond
	}
	defer con.Close()
	return true, latency

}

func printColorText(connected bool, latency time.Duration) {
	fmt.Printf("\033[2J") // Clear screen
	fmt.Printf("\033[H")  // Move cursor to home

	if connected {
		success.Printf("Connected to %v", url)
	} else {
		danger.Printf("Not connected to %v", url)
	}
	warning.Print(" : ")
	if latency < 1*time.Millisecond {
		danger.Print(latency)
	} else if latency < 30*time.Millisecond {
		success.Print(latency)
	} else if latency < 100*time.Millisecond {
		warning.Print(latency)
	} else {
		danger.Print(latency)
	}

	success.Print("\n")
}

func printLatencyGraph(latency time.Duration) {
	latencyRing.Enqueue(latency)
	latencyArray := make([]float64, latencyRing.Length)
	line := danger.Sprint("  -1  ")

	for i := 0; i < latencyRing.Length; i += 1 {
		number := float64(latencyRing.Get(i).Milliseconds())
		latencyArray[i] = number
		if number < 0 {
			line += danger.Sprint("V")
		} else {
			line += " "
		}
	}
	graph := asciigraph.Plot(latencyArray, asciigraph.Height(20), asciigraph.Precision(0))

	fmt.Println(graph)
	fmt.Println(line)

}
