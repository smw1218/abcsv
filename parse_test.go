package abcsv

import (
	"strings"
	"testing"
)

var testOutput = `
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
 100%      6 (longest request)`

func TestParse(t *testing.T) {
	reader := strings.NewReader(testOutput)
	results := ParseAB(reader)
	if results.Hostname != "localhost" {
		t.Errorf("Failed host parse expected \"localhost\" got \"%v\"", results.Hostname)
	}
	if results.Concurrency != 2 {
		t.Errorf("Failed concurrency parse expected 2 got %v", results.Concurrency)
	}
	if results.AverageResponseTime != 2.178 {
		t.Errorf("Failed AverageResponseTime parse expected 2.178 got %v", results.AverageResponseTime)
	}
}
