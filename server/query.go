package server

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/mzki/feserver/src"
)

// TODO: The name of URL query parameters are
// configurable by external file?

const (
	// query parameters for getQuestion.
	QueryYear   = "year"
	QuerySeason = "season"
	QueryNo     = "no"
)

func parseGetQuestionQuery(v url.Values, source Source) (src.Query, error) {
	q := src.Query{}
	q.Year = parseIntParam(v, QueryYear, q.Year)
	q.No = parseIntParam(v, QueryNo, q.No)
	if s := v.Get(QuerySeason); s != "" {
		q.Season = s
	}
	if err := source.Source.Validates(q); err != nil {
		return q, fmt.Errorf("invalid query form (%s), %v", v.Encode(), err)
	}
	return q, nil
}

const (
	// query parameters for getRandom.
	QueryMaxYear     = "max_year"
	QueryMinYear     = "min_year"
	QueryMaxNo       = "max_no"
	QueryMinNo       = "min_no"
	QuerySeasonRange = "season"
)

func parseGetRandomQuery(v url.Values, source Source) (src.QueryRange, error) {
	qr := source.Source.QueryRange
	qr.MaxYear = parseIntParam(v, QueryMaxYear, qr.MaxYear)
	qr.MinYear = parseIntParam(v, QueryMinYear, qr.MinYear)
	qr.MaxNo = parseIntParam(v, QueryMaxNo, qr.MaxNo)
	qr.MinNo = parseIntParam(v, QueryMinNo, qr.MinNo)
	if s := v.Get(QuerySeasonRange); s != "" {
		qr.Season = s
	}
	if err := source.Source.ValidatesRange(qr); err != nil {
		return qr, fmt.Errorf("invalid query form (%s), %v", v.Encode(), err)
	}
	return qr, nil
}

func parseIntParam(v url.Values, key string, _default int) int {
	if param := v.Get(key); param != "" {
		if i, err := strconv.Atoi(param); err == nil {
			return i
		}
	}
	return _default
}
