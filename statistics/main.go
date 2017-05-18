//go:generate go run main.go
//
// +build stat

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mzki/feserver/src"
)

const File = "./statistics.csv"

func main() {
	fp, err := os.Create(File)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	writeStatistics(fp)
}

func writeStatistics(writer io.Writer) {
	fmt.Fprintf(writer, "year, season, imageCount, rate, techRate, manageRate, strategyRate, errCount\n")

	totalCount := 0
	totalTech := 0
	totalMana := 0
	totalStrat := 0
	cycle := 0

	for y := src.MinYear; y <= src.MaxYear; y++ {
		for _, s := range []string{src.SeasonSpring, src.SeasonAutumn} {
			count, techC, manaC, stratC, errCount := countHasImage(y, s)

			totalCount += count
			totalTech += techC
			totalMana += manaC
			totalStrat += stratC
			cycle += 1

			fmt.Fprintf(writer, "%d, %s, %d/%d, %.2f, %.2f, %.2f, %.2f, %d\n",
				y, s, count, src.MaxNo, rate(count), rate(techC), rate(manaC), rate(stratC), errCount)
		}
	}

	totalQ := cycle * src.MaxNo
	fmt.Fprintf(writer, "# total, imageCount, rate, techRate, manageRate, strategyRate\n")
	fmt.Fprintf(writer, "# - , %d/%d, %.2f, %.2f, %.2f, %.2f\n",
		totalCount, totalQ,
		float64(totalCount)/float64(totalQ),
		float64(totalTech)/float64(totalQ),
		float64(totalMana)/float64(totalQ),
		float64(totalStrat)/float64(totalQ),
	)
}

func rate(a int) float64 {
	return float64(a) / float64(src.MaxNo)
}

const (
	TechMax  = 50
	ManaMax  = 60
	StratMax = 80
)

var getter = src.NewGetter(src.FE, src.LeastIntervalTime)

func countHasImage(year int, season string) (count, techC, manaC, stratC, errCount int) {
	ctx := context.Background()
	for q := src.MinNo; q <= src.MaxNo; q++ {
		res, err := getter.Get(ctx, src.Query{
			Year:   year,
			Season: season,
			No:     q,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			errCount += 1
			continue
		}
		if res.HasImage {
			count++
			switch {
			case q <= TechMax:
				techC += 1
			case q <= ManaMax:
				manaC += 1
			case q <= StratMax:
				manaC += 1
			}
		}
	}
	return
}
