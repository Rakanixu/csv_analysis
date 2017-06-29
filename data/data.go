package data

import (
	"fmt"
	"log"
	"strings"
	"time"

	"encoding/json"

	"github.com/Rakanixu/csv_analysis/db"
	"github.com/nu7hatch/gouuid"
)

const (
	PADDING = 45
)

// Data ...
type Data struct {
	date            time.Time
	name            string
	Records         map[string]*Record
	NumTotalColumns int64
}

// Record ...
type Record struct {
	description string
	count       int64
	percentage  float64
}

// NewData ...
func NewData(name string) *Data {
	return &Data{
		name:    name,
		Records: make(map[string]*Record),
	}
}

// NewRecord ...
func NewRecord(description string) *Record {
	return &Record{
		description: description,
	}
}

// AddRecord ...
func (d *Data) AddRecord(record *Record) {
	if d.Records[record.description] == nil {
		d.Records[record.description] = record
	}
	d.Records[record.description].count++
}

// SetTotal ...
func (d *Data) SetTotal(total int64) {
	d.NumTotalColumns = total
}

// Date ...
func (d *Data) Date() {
	var err error

	// Genetare Time from csv file name pattern
	// xxxxxxx-1c2b447a65025fc41ee91b794aa2305a-11042016.csv
	n := strings.Split(d.name, "-")
	s := strings.Replace(n[len(n)-1], ".csv", "", -1)

	s = fmt.Sprintf("%s%s/%s%s/%s%s%s%s",
		string(s[2]), string(s[3]), // MM
		string(s[0]), string(s[1]), // DD
		string(s[4]), string(s[5]), string(s[6]), string(s[7]), // YYYYY
	)

	d.date, err = time.Parse("01/02/2006", s)
	if err != nil {
		log.Println("PARSE CSV TIME ERR", err)
	}
}

// Info ...
func (d *Data) Info() {
	for _, v := range d.Records {
		v.percentage = float64(v.count) / float64(d.NumTotalColumns) * 100
	}
}

// Print ...
func (d *Data) Print() {
	d.Info()

	fmt.Println("\n-------------------------------------------------------------------")
	fmt.Println("DATE:                                       ", d.date)
	fmt.Println("CSV:                                        ", d.name)
	for _, v := range d.Records {
		fmt.Print(v.description)
		if len(v.description)-PADDING < 0 {
			fmt.Printf("\033[%dC", len(v.description)-PADDING)
		} else {
			fmt.Println("")
			fmt.Printf("\033[%dC", PADDING)
		}
		fmt.Printf("%2.3f", v.percentage)
		fmt.Print("%   ")
		fmt.Print(v.count)
		fmt.Println("")
	}
}

// Dump ...
func (d *Data) Dump() error {
	type info struct {
		ID          string  `json:"id"`
		Count       int64   `json:"count"`
		Total       int64   `json:"total"`
		Description string  `json:"description"`
		Percentage  float64 `json:"percentage"`
		Date        string  `json:"date"`
		CsvFile     string  `json:"csv_file"`
	}

	for _, v := range d.Records {
		u, err := uuid.NewV4()
		if err != nil {
			log.Println("ERROR GENRATING UUID", err)
		}

		i := info{
			ID:          u.String(),
			Count:       v.count,
			Total:       d.NumTotalColumns,
			Description: v.description,
			Percentage:  v.percentage,
			Date:        d.date.String(),
			CsvFile:     d.name,
		}

		/*		db.BulkIndex(d.date.String(), )*/

		b, err := json.Marshal(i)
		if err != nil {
			log.Println("ERROR MARSHALLING", err)
		}

		if err := db.Index(u.String(), string(b)); err != nil {
			log.Println("ERROR INDEXING", err)
		}
	}

	return nil
}

// DataSlice ...
type DataSlice []*Data

// Len ...
func (d DataSlice) Len() int {
	return len(d)
}

// Less ...
func (d DataSlice) Less(i, j int) bool {
	return d[i].date.After(d[j].date)
}

// Swap ...
func (d DataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
