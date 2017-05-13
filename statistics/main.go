// +build stat

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mzki/feserver/src"
)

func main() {
	totalCount := 0
	cycle := 0
	for y := src.MinYear; y <= src.MaxYear; y++ {
		for _, s := range []string{src.SeasonSpring, src.SeasonAutumn} {
			count, errCount := countHasImage(y, s)
			totalCount += count
			cycle += 1
			fmt.Fprintf(os.Stdout, "year:%d, season:%s, imageCount:%d/%d, rate:%.1f, errCount:%d\n",
				y, s, count, src.MaxNO, float64(count)/float64(src.MaxNO), errCount)
		}
	}
	fmt.Fprintf(os.Stdout, "# total, imageCount:%d/%d, rate:%.2f\n",
		totalCount, cycle*src.MaxNO, float64(totalCount)/float64(cycle*src.MaxNO))
}

var getter = src.NewGetter(src.LeastIntervalTime)

func countHasImage(year int, season string) (int, int) {
	ctx := context.Background()
	count := 0
	errCount := 0
	for q := src.MinNO; q <= src.MaxNO; q++ {
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
		}
	}
	return count, errCount
}
