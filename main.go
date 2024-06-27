package main

import (
	"fmt"
	duration_circular_buffer "internet-checker/circular_buffer"
	"net"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
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
		return false, -1
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
	if latency < 30*time.Millisecond {
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

	lines := make([]string, 4)
	for i := 0; i < latencyRing.Length; i += 1 {
		latency := latencyRing.Get(i)
		if latency == -1 {
			lines[0] += " "
			lines[1] += " "
			lines[2] += " "
			lines[3] += danger.Sprint("V")
		} else if latency < 30*time.Millisecond {
			lines[0] += " "
			lines[1] += " "
			lines[2] += success.Sprint("-")
			lines[3] += " "
		} else if latency < 100*time.Millisecond {
			lines[0] += " "
			lines[1] += warning.Sprint("^")
			lines[2] += warning.Sprint("-")
			lines[3] += " "
		} else {
			lines[0] += danger.Sprint("^")
			lines[1] += warning.Sprint("|")
			lines[2] += warning.Sprint("-")
			lines[3] += " "
		}

	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
