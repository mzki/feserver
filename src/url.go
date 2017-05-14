package src

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"text/template"
	"time"
)

var targetURLTmpl = template.Must(template.New("targetURL").Parse(
	`http://www.fe-siken.com/kakomon/{{.Year}}_{{.Season}}/q{{.No}}.html`))

const (
	// These represents the minimum and maximum query range.

	// Examination Year range
	MinYear = 13
	MaxYear = 29

	// Examination Season selections
	SeasonSpring = "haru"
	SeasonAutumn = "aki"
	SeasonAll    = "all" // used for QueryRange only.

	// Question No.
	MinNo = 1
	MaxNo = 80
)

var (
	yearRange   = [...]int{MinYear, MaxYear}
	seasonRange = [...]string{SeasonSpring, SeasonAutumn}
	noRange     = [...]int{MinNo, MaxNo}
)

// Query is a query for source URL.
// Its fields specifies which question is searched for.
type Query struct {
	Year   int
	Season string
	No     int
}

// check whether query has correct value range?
// nil error means query is valid.
func (q Query) Validates() error {
	if y := q.Year; y < MinYear || y > MaxYear {
		return fmt.Errorf("year must be in [%d:%d], but %d", MinYear, MaxYear, y)
	}
	if s := q.Season; s != SeasonSpring && s != SeasonAutumn {
		return fmt.Errorf("Season must be either %s or %s, but %s", SeasonSpring, SeasonAutumn, s)
	}
	if n := q.No; n < MinNo || n > MaxNo {
		return fmt.Errorf("Question No. must be in [%d:%d], but %d", MinNo, MaxNo, n)
	}
	return nil
}

var DefaultQueryRange = QueryRange{
	MaxYear: MaxYear, MinYear: MinYear,
	MaxNo: MaxNo, MinNo: MinNo,
	Season: SeasonAll,
}

// QueryRange represents query range for randomly selected.
type QueryRange struct {
	MaxYear, MinYear int
	MaxNo, MinNo     int
	Season           string // SeasonSpring | SeasonAutumn | SeasonAll
}

// check whether QueryRange is in correct range?
func (qr QueryRange) Validates() error {
	if err := (Query{qr.MaxYear, SeasonSpring, qr.MaxNo}).Validates(); err != nil {
		return err
	}
	if err := (Query{qr.MinYear, SeasonSpring, qr.MinNo}).Validates(); err != nil {
		return err
	}

	// check the relation for min and max.
	if qr.MaxYear < qr.MinYear {
		return fmt.Errorf("src.QueryRange: MaxYear must be larger then MinYear but Max: %d, Min: %d", qr.MaxYear, qr.MinYear)
	}
	if qr.MaxNo < qr.MinNo {
		return fmt.Errorf("src.QueryRange: MaxNo must be larger then MinNo but Max: %d, Min: %d", qr.MaxNo, qr.MinNo)
	}
	// check season
	switch s := qr.Season; s {
	case SeasonAll, SeasonSpring, SeasonAutumn:
		return nil
	default:
		return fmt.Errorf("src.QueryRange: Season must be %s, %s or %s", SeasonSpring, SeasonAutumn, SeasonAll)
	}
}

func (qr QueryRange) season(r *rand.Rand) string {
	switch s := qr.Season; s {
	case SeasonSpring, SeasonAutumn:
		return s
	case SeasonAll:
		return seasonRange[random.Intn(len(seasonRange))]
	default:
		panic("src.QueryRange: unknown season " + s)
		return ""
	}
}

var randMutex = new(sync.Mutex)

// package global random state. under mutex.
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// generates random query.
func randomQuery(pqr *QueryRange) Query {
	qrange := *pqr

	beforeAutumn := !autumnPublished()
	randMutex.Lock()
	// autumn examination in the latest year is unpublisded,
	// shirnk the max year range.
	season := qrange.season(random)
	if beforeAutumn && season == SeasonAutumn && qrange.MaxYear == MaxYear {
		qrange.MaxYear -= 1
	}
	year := random.Intn(qrange.MaxYear-qrange.MinYear+1) + qrange.MinYear
	no := random.Intn(qrange.MaxNo-qrange.MinNo+1) + qrange.MinNo
	randMutex.Unlock()
	return Query{year, season, no}
}

// indicates the day in which the autumn examination is published.
const autumnPublishedMonth = time.November

func autumnPublished() bool {
	now := time.Now()
	publishedDay := time.Date(now.Year(), autumnPublishedMonth, 0, 0, 0, 0, 0, time.UTC)
	return now.After(publishedDay)
}

// generate target URL with randomized query.
// randomized range is defined by QueryRange.
// use default range if QueryRange is nil
func RandomURL(qr *QueryRange) string {
	if qr == nil {
		qr = &DefaultQueryRange
	}
	q := randomQuery(qr)
	return GenerateURL(q)
}

// generate source URL with given query.
// the query must have valid fields range
// which can be validated by (Query).Validates().
func GenerateURL(q Query) string {
	if err := q.Validates(); err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := targetURLTmpl.Execute(buf, &q); err != nil {
		panic(err) // TODO
	}
	return buf.String()
}
