package main

/*
Usage: go run main.go -c "Error Code" -dimension_key Device -dimension_val Chromecast
*/

// TODO: Get the bigger size of all CSV to process so we do not waste memory
import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sync"

	"github.com/Rakanixu/csv_analysis/data"
)

const (
	ERR_CODE = "Error Code"
)

var results []*data.Data
var (
	key, value, dimension *string
	maxGoroutines         *int
	wg                    sync.WaitGroup
)

func main() {
	p := flag.String("p", "", "Path to CSV files")
	maxGoroutines = flag.Int("t", 8, "Number of parallel gourotines")
	dimension = flag.String("c", ERR_CODE, "Column / dimension to apply aggregation")
	key = flag.String("dimension_key", "Device", "Dimension key name")
	value = flag.String("dimension_val", "Chromecast", "Dimension value")
	flag.Parse()

	if *dimension == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	analyzeCSVs(getCSVFiles(*p))

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

func analyzeCSVs(paths []string) {
	// Channel buffer equal to number of gourotines
	blocker := make(chan struct{}, *maxGoroutines)

	for _, v := range paths {
		// Fill blocker channel
		blocker <- struct{}{}
		wg.Add(1)

		go func(path string) {
			s, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer s.Close()

			fi, err := s.Stat()
			log.Println("Size", fi.Name(), fi.Size())

			// Increase size if CSV file is > 500MB
			b := make([]byte, fi.Size())
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

			analyzeCSV(path, records[1:], i, j)
			// Read from blocker channel to allow next iteration
			<-blocker
			wg.Done()
		}(v)
	}

	wg.Wait()
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
			if len(r) > 1 && len(r) > aggDimensionIndex && len(r) > filterIndex && !(f && r[filterIndex] != *value) {
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
