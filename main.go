package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func read(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var targetDates []time.Time

	for i, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}

		ts, err := time.Parse("Jan 2 3:04:05 PM MST 2006", line+" 2025")
		if err != nil {
			return fmt.Errorf("error parsing record %d: %w", i, err)
		}

		if ts.Month() == time.January || ts.Month() == time.February {
			if time.Since(ts) < 14*24*time.Hour {
				targetDates = append(targetDates, ts)
			}
		}
	}

	r := csv.NewReader(f)
	r.LazyQuotes = true

	var d time.Duration
	for i := 0; i < len(targetDates)-1; i++ {
		j := i + 1

		gap := targetDates[j].Sub(targetDates[i])
		fmt.Println(gap)
		if gap > d {
			d = gap
		}
	}

	fmt.Println("largest gap", d)

	return nil
}

// cat elonmusk.csv |  grep -E 'Jan|Feb' | awk -F, '{print $(NF-1) $NF}' | sed 's/"//g' | grep -E '^(Jan|Feb)' > timestamps.txt
func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
	}

	if err := read(flag.Arg(0)); err != nil {
		log.Fatal(err)
	}
}
