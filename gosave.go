// gosave takes CSV input and calculates monthly savings
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type Savings struct {
	Date    DateTime
	Balance float64
}

var netIncome int
var totalSaving float64

func main() {
	var err error
	var csvFile string

	flag.StringVar(&csvFile, "file", "", "CSV file name")
	flag.IntVar(&netIncome, "income", 0, "Net income (after tax)")

	flag.Parse()

	if csvFile == "" {
		fmt.Printf("CSV file must be specified\n")
		flag.PrintDefaults()
		return
	}

	savingsFile, err := os.OpenFile(csvFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	defer savingsFile.Close()

	savings := []Savings{}

	if err := gocsv.UnmarshalFile(savingsFile, &savings); err != nil { // Load savings from file
		fmt.Printf("Error %s\n", err)
		return
	}

	// find direction of travel
	var ascending bool
	var time0 time.Time
	for _, row := range savings {
		if time0.IsZero() {
			time0 = row.Date.Time
		} else if row.Date.After(time0) {
			ascending = true
			break
		} else if row.Date.Equal(time0) {
			continue
		} else {
			break
		}
	}

	var balanceEnd, balanceStart float64
	var month time.Time
	var months int
	for i, row := range savings {
		//fmt.Printf("%+v\n", row)
		if month.IsZero() {
			month = row.Date.Time
			if ascending {
				balanceStart = row.Balance
			} else {
				balanceEnd = row.Balance
			}
		} else if month.Month() != row.Date.Month() {
			months++
			//fmt.Printf("%10s %d, saved $%10.2f, balanceStart $%10.2f, balanceEnd $%10.2f\n", month.Month(), month.Year(), balanceEnd-balanceStart, balanceStart, balanceEnd)
			PrintSaving(month, balanceEnd-balanceStart)
			month = row.Date.Time
			if ascending {
				balanceStart = row.Balance
			} else {
				balanceEnd = row.Balance
			}
		}
		if ascending {
			balanceEnd = row.Balance
		} else {
			balanceStart = row.Balance
		}

		// last row
		if i+1 == len(savings) {
			PrintSaving(month, balanceEnd-balanceStart)
		}
	}
	var pc float64
	if netIncome > 0 {
		adjustedIncome := float64(netIncome) * float64(months) / 12
		pc = (totalSaving / float64(adjustedIncome)) * 100
	}
	fmt.Printf("\nTotal Saved: $%10.2f %7.2f%%\n", totalSaving, pc)
}

func PrintSaving(month time.Time, saving float64) {
	totalSaving += saving
	fmt.Printf("%10s %d, saved $%10.2f\n", month.Month(), month.Year(), saving)
}

type DateTime struct {
	time.Time
}

func (d *DateTime) UnmarshalCSV(csv string) (err error) {
	//s.Date, err = time.Parse("02/01/2006", csv)
	d.Time, _ = time.Parse("02/01/2006", csv)
	/*
		if err != nil {
			return err
		}
	*/
	return nil
}

/*
func (s *Savings) MarshalCSV() (string, error) {
	return s.Date.Format("02/01/2006"), nil
}

func (s *Savings) String() string {
	return s.String() // Redundant, just for example
}
*/
