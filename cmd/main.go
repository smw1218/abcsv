package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/smw1218/abcsv"
)

func main() {
	var name string
	var addHeader bool
	var onlyHeader bool
	flag.StringVar(&name, "n", "", "optional name of the test run")
	flag.BoolVar(&addHeader, "h", false, "print csv header before results")
	flag.BoolVar(&onlyHeader, "H", false, "print csv header and exit")
	flag.Parse()

	if onlyHeader {
		fmt.Println(abcsv.Columns())
		return
	}

	res := abcsv.ParseAB(os.Stdin)
	if addHeader {
		fmt.Println(abcsv.Columns())
	}
	fmt.Println(res.Csv(name))
}
