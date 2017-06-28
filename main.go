package main

/*
Usage: go run main.go -c "Error Code" -device_key Device -device_val Chromecast
*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Rakanixu/csv_analysis/data"
)

var results []*data.Data
var (
	deviceKey, deviceName *string
)

func main() {
	bs := flag.Int64("m", 800000000, "CSV file size (bytes) default to 800MB")
	col := flag.String("c", "Error Code", "Column / dimension to apply aggregation")
	deviceKey = flag.String("device_key", "Device", "Device key name")
	deviceName = flag.String("device_val", "Chromecast", "Device value name")
	flag.Parse()

	if *col == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	analyzeCSVs(getCSVFiles(), *bs, *col)

	sort.Sort(data.DataSlice(results))
	for _, v := range results {
		v.Print()
	}
}

func getCSVFiles() []string {
	files, err := filepath.Glob("*.csv")
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func analyzeCSVs(paths []string, size int64, colName string) {
	for _, v := range paths {
		s, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}

		// Increase size if CSV file is > 500MB
		b := make([]byte, size)
		count, err := s.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		records := strings.Split(string(b[:count]), "\n")
		columms := strings.Split(records[0], ",")

		i := -1
		j := -1
		for k, v := range columms {
			// CSV can contain CDN or "CDN"
			switch trimDoubleQuote(v) {
			case colName:
				i = k
			case *deviceKey:
				j = k
			}
		}

		analyzeCSV(v, records[1:], i, j)
	}
}

func analyzeCSV(name string, csv []string, index, deviceIndex int) {
	var n int64
	d := data.NewData(name)
	f := false

	if deviceIndex > 0 && len(*deviceName) > 0 {
		f = true
	}

	for _, v := range csv {
		r := strings.Split(v, ",")

		// Don't push records which type is different to the one set on flags
		if len(r) > 1 && !(f && r[deviceIndex] != *deviceName) {
			n++
			d.AddRecord(data.NewRecord(fmt.Sprintf("%s %s", r[index], r[1])))
		}
	}
	d.SetTotal(n)
	d.Date()

	results = append(results, d)
}

func trimDoubleQuote(s string) string {
	return strings.Replace(s, `"`, "", -1)
}
