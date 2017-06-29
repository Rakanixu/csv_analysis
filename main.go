package main

/*
Usage: go run main.go -c "Error Code" -dimension_key Device -dimension_val Chromecast
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

const (
	ERR_CODE = "Error Code"
)

var results []*data.Data
var (
	key, value, dimension *string
)

func main() {
	bs := flag.Int64("m", 800000000, "CSV file size (bytes) default to 800MB")
	p := flag.String("p", "", "path to CSV files")
	dimension = flag.String("c", ERR_CODE, "Column / dimension to apply aggregation")
	key = flag.String("dimension_key", "Device", "Dimension key name")
	value = flag.String("dimension_val", "Chromecast", "Dimension value")
	flag.Parse()

	if *dimension == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	analyzeCSVs(getCSVFiles(*p), *bs)

	sort.Sort(data.DataSlice(results))
	for _, v := range results {
		v.Print()
	}
}

func getCSVFiles(path string) []string {
	if path != "" {
		path = fmt.Sprintf("%s/*.csv", path)
	} else {
		path = "*.csv"
	}

	files, err := filepath.Glob(path)
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func analyzeCSVs(paths []string, size int64) {
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
			case *dimension:
				i = k
			case *key:
				j = k
			}
		}

		analyzeCSV(v, records[1:], i, j)
	}
}

func analyzeCSV(name string, csv []string, aggDimensionIndex, filterIndex int) {
	if aggDimensionIndex >= 0 && filterIndex >= 0 {
		var n int64
		d := data.NewData(name)
		f := false

		if filterIndex > 0 && len(*value) > 0 {
			f = true
		}

		for _, v := range csv {
			r := strings.Split(v, ",")

			// Don't push records which type is different to the one set on flags
			if len(r) > 1 && len(r) > aggDimensionIndex && !(f && r[filterIndex] != *value) {
				des := r[aggDimensionIndex]
				n++

				if *dimension == ERR_CODE {
					// HARDCODED: specific case for a CSV specific pattern
					des = fmt.Sprintf("%s %s", r[aggDimensionIndex], r[1])
				}
				d.AddRecord(data.NewRecord(des))
			}
		}
		d.SetTotal(n)
		d.Date()

		results = append(results, d)
	}
}

func trimDoubleQuote(s string) string {
	return strings.Replace(s, `"`, "", -1)
}
