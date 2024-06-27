package main

import (
	"fmt"
	"net"
	duration_circular_buffer "ping-checker/duration-circular_buffer"
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

	lines := make([]string, 6)
	for i := 0; i < latencyRing.Length; i += 1 {
		latency := latencyRing.Get(i)
		if latency == -1 {
			lines[0] += " "
			lines[1] += " "
			lines[2] += " "
			lines[3] += " "
			lines[4] += " "
			lines[5] += danger.Sprint("V")
		} else if latency < 10*time.Millisecond {
			lines[0] += " "
			lines[1] += " "
			lines[2] += " "
			lines[3] += " "
			lines[4] += success.Sprint("▮")
			lines[5] += " "
		} else if latency < 30*time.Millisecond {
			lines[0] += " "
			lines[1] += " "
			lines[2] += " "
			lines[3] += success.Sprint("▮")
			lines[4] += success.Sprint("▮")
			lines[5] += " "
		} else if latency < 70*time.Millisecond {
			lines[0] += " "
			lines[1] += " "
			lines[2] += warning.Sprint("▮")
			lines[3] += warning.Sprint("▮")
			lines[4] += warning.Sprint("▮")
			lines[5] += " "
		} else if latency < 100*time.Millisecond {
			lines[0] += " "
			lines[1] += danger.Sprint("▮")
			lines[2] += danger.Sprint("▮")
			lines[3] += danger.Sprint("▮")
			lines[4] += danger.Sprint("▮")
			lines[5] += " "
		} else {
			lines[0] += danger.Sprint("▮")
			lines[1] += danger.Sprint("▮")
			lines[2] += danger.Sprint("▮")
			lines[3] += danger.Sprint("▮")
			lines[4] += danger.Sprint("▮")
			lines[5] += " "
		}

	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
