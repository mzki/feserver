//go:generate gopherjs build -o fesrc.js -m

package main

import (
	"context"
	"log"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/mzki/feserver/src"
)

func TimeoutCtx(timeout uint64) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if timeout > 0 {
		return context.WithTimeout(ctx, time.Duration(timeout))
	}
	return ctx, context.CancelFunc(func() {})
}

// because http package is not supported by gopherjs,
// Deprecated
func GetRandom(timeout uint64) *src.Response {
	ctx, cancel := TimeoutCtx(timeout)
	defer cancel()
	res, err := src.GetRandom(ctx, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &res
}

// because http package is not supported by gopherjs,
// Deprecated
func Get(year int, season string, no int, timeout uint64) *src.Response {
	ctx, cancel := TimeoutCtx(timeout)
	defer cancel()
	res, err := src.Get(ctx, src.Query{
		Year:   year,
		Season: season,
		No:     no,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return &res
}

func ParseHTML(html string) *src.Response {
	res, err := src.ParseHTML(html)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &res
}

func main() {
	js.Global.Set("fesrc", map[string]interface{}{
		"ParseHTML": ParseHTML,
		"RandomURL": src.RandomURL,
	})
}
