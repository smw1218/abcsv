package abcsv

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	firstLine          = regexp.MustCompile(`^This is ApacheBench, Version ([0-9.]+) <\$Revision: (\d+) \$>`)
	kvLine             = regexp.MustCompile(`([\w ]+):\s+(\S+.*)`)
	leadingNumber      = regexp.MustCompile(`([0-9.]+)(.*)`)
	connectionTimesRow = regexp.MustCompile(`([0-9.]+)\s+([0-9.]+)\s+([0-9.]+)\s+([0-9.]+)\s+([0-9.]+)\s*`)
	percentLine        = regexp.MustCompile(`\s*(\d+)%\s+(.*)`)
)

func ParseAB(r io.Reader) *Results {
	scanner := bufio.NewScanner(r)
	res := &Results{
		ConnectionTimes: &ConnectionTimes{},
		NTiles:          make(map[int]float64),
	}
	for scanner.Scan() {
		line := scanner.Text()
		if vals := firstLine.FindStringSubmatch(line); len(vals) == 3 {
			res.Version = parseFloat(vals[1])
			res.Revision = parseInt(vals[2])
		} else if vals = kvLine.FindStringSubmatch(line); len(vals) == 3 {
			//fmt.Printf("Matched %d k: %v v: %v\n", len(vals), vals[1], vals[2])
			setMatched(res, vals[1], vals[2])
		} else if vals = percentLine.FindStringSubmatch(line); len(vals) == 3 {
			//fmt.Printf("Matched percent %d k: %v v: %v\n", len(vals), vals[1], vals[2])
			setNTile(res, vals[1], vals[2])
		} else {
			//fmt.Printf("No match: %v\n", line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	return res
}

func setNTile(res *Results, percentS string, timeS string) {
	percent := parseInt(percentS)
	if percent > 0 {
		res.NTiles[percent] = parseFloat(timeS)
	}
}

func setMatched(res *Results, key string, value string) {
	switch key {
	case "Server Software":
		res.Server = value
	case "Server Hostname":
		res.Hostname = value
	case "Server Port":
		res.Port = parseInt(value)
	case "Document Path":
		res.Path = value
	case "Document Length":
		res.BodySize = parseInt(value)
	case "Concurrency Level":
		res.Concurrency = parseInt(value)
	case "Time taken for tests":
		res.TestTime = parseFloat(value)
	case "Complete requests":
		res.CompletedRequests = parseInt(value)
	case "Failed requests":
		res.FailedRequests = parseInt(value)
	case "Total transferred":
		res.TotalSize = parseInt(value)
	case "HTML transferred":
		res.BodySizeTotal = parseInt(value)
	case "Requests per second":
		res.Throughput = parseFloat(value)
	case "Time per request":
		setTimePerRequest(value, res)
	case "Transfer rate":
		res.TransferRate = parseFloat(value)
	case "Connect":
		res.ConnectionTimes.Connect = parseConnection(value)
	case "Processing":
		res.ConnectionTimes.Processing = parseConnection(value)
	case "Waiting":
		res.ConnectionTimes.Waiting = parseConnection(value)
	case "Total":
		res.ConnectionTimes.Total = parseConnection(value)
	}
}

func setTimePerRequest(value string, res *Results) {
	vals := leadingNumber.FindStringSubmatch(value)
	if len(vals) != 3 {
		log.Printf("Fail submatch!!!")
		return
	}
	if vals[2] == " [ms] (mean)" {
		res.AverageResponseTime = parseFloat(value)
	}
}

func parseConnection(value string) *ConnectionTimeStats {
	vals := connectionTimesRow.FindStringSubmatch(value)
	if len(vals) != 6 {
		return nil
	}
	cts := &ConnectionTimeStats{
		Min:    parseFloat(vals[1]),
		Mean:   parseFloat(vals[2]),
		Std:    parseFloat(vals[3]),
		Median: parseFloat(vals[4]),
		Max:    parseFloat(vals[5]),
	}
	return cts
}

/*
This is ApacheBench, Version 2.3 <$Revision: 1826891 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient).....done


Server Software:
Server Hostname:        localhost
Server Port:            6060

Document Path:          /
Document Length:        6728 bytes

Concurrency Level:      2
Time taken for tests:   0.011 seconds
Complete requests:      10
Failed requests:        0
Total transferred:      68240 bytes
HTML transferred:       67280 bytes
Requests per second:    918.11 [#/sec] (mean)
Time per request:       2.178 [ms] (mean)
Time per request:       1.089 [ms] (mean, across all concurrent requests)
Transfer rate:          6118.31 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       0
Processing:     0    2   1.8      1       6
Waiting:        0    2   1.8      1       6
Total:          0    2   1.8      1       6

Percentage of the requests served within a certain time (ms)
  50%      1
  66%      2
  75%      3
  80%      4
  90%      6
  95%      6
  98%      6
  99%      6
 100%      6 (longest request)
*/

func parseInt(value string) int {
	i, err := strconv.Atoi(scrubString(value))
	if err != nil {
		log.Printf("Failed parse: %v", err)
	}
	return i
}

func parseFloat(value string) float64 {
	f, err := strconv.ParseFloat(scrubString(value), 64)
	if err != nil {
		log.Printf("Failed parse: %v", err)
	}
	return f
}

func scrubString(value string) string {
	vals := leadingNumber.FindStringSubmatch(value)
	if len(vals) != 3 {
		log.Printf("Fail submatch!!!")
		return "0"
	}
	return vals[1]
}

type Results struct {
	Version             float64
	Revision            int
	Server              string
	Hostname            string
	Port                int
	Path                string
	BodySize            int // bytes
	Concurrency         int
	TestTime            float64 // s
	CompletedRequests   int
	FailedRequests      int
	TotalSize           int              // bytes
	BodySizeTotal       int              // bytes
	Throughput          float64          // req/s
	AverageResponseTime float64          // ms
	TransferRate        float64          // kb/s
	ConnectionTimes     *ConnectionTimes // ms
	NTiles              map[int]float64  // ms
}

func Columns() string {
	return "Name,Server,Hostname,Port,Path,Concurrency,Throughput,Avg. Latency,Duration,Successful,Failed,Max. latency,50% Latency,90% Latency,95% Latency,98% Latency,99% Latency,Avg. Recv. Bandwidth"
}

func (r *Results) Csv(name string) string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		name,
		r.Server,
		r.Hostname,
		r.Port,
		r.Path,
		r.Concurrency,
		r.Throughput,
		r.AverageResponseTime,
		r.TestTime,
		r.CompletedRequests,
		r.FailedRequests,
		r.NTiles[100],
		r.NTiles[50],
		r.NTiles[90],
		r.NTiles[95],
		r.NTiles[98],
		r.NTiles[99],
		r.TransferRate,
	)
}

type ConnectionTimes struct {
	Connect    *ConnectionTimeStats
	Processing *ConnectionTimeStats
	Waiting    *ConnectionTimeStats
	Total      *ConnectionTimeStats
}

type ConnectionTimeStats struct {
	Min    float64
	Mean   float64
	Std    float64
	Median float64
	Max    float64
}

type NTile struct {
	Percent      int
	ResponseTime float64
}
