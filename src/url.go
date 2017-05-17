package src

import (
	"bytes"
	"math/rand"
	"sync"
	"text/template"
	"time"
)

type urlGenerator struct {
	src  Source
	tmpl *template.Template
}

func newURLGenerator(s Source) *urlGenerator {
	if err := s.ValidatesSelf(); err != nil {
		panic(err)
	}
	return &urlGenerator{
		src:  s,
		tmpl: template.Must(template.New(s.URL).Parse(s.URL)),
	}
}

// generate source URL with given query.
func (url *urlGenerator) Generate(q Query) (string, error) {
	if err := url.src.Validates(q); err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := url.tmpl.Execute(buf, &q); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// generate randomized source URL with query range.
func (url *urlGenerator) Random(qr QueryRange) (string, error) {
	if qr == MaxQueryRange {
		qr = url.MaxQueryRange()
	}
	if err := url.src.ValidatesRange(qr); err != nil {
		return "", err
	}
	return url.Generate(randomQuery(qr))
}

// return maximum range of query for the url's source.
func (url *urlGenerator) MaxQueryRange() QueryRange {
	return url.src.QueryRange
}

// These represents the minimum and maximum query range
// for the F.E. examination.
const (
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

// MaxQueryRange indicates the maximum range of query in the source.
var MaxQueryRange = QueryRange{-1, -1, -1, -1, SeasonAll}

// QueryRange represents query range for randomly selected.
type QueryRange struct {
	MaxYear, MinYear int
	MaxNo, MinNo     int
	Season           string // SeasonSpring | SeasonAutumn | SeasonAll
}

func (qr QueryRange) season(r *rand.Rand) string {
	switch s := qr.Season; s {
	case SeasonSpring, SeasonAutumn:
		return s
	case SeasonAll:
		return seasonRange[r.Intn(len(seasonRange))]
	default:
		panic("src.QueryRange: unknown season " + s)
		return ""
	}
}

var randMutex = new(sync.Mutex)

// package global random state. under mutex.
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// generates random query.
func randomQuery(qrange QueryRange) Query {
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

var _FE_URL = newURLGenerator(FE)

// generate random target URL.
// randomized range is defined by given QueryRange.
// use default range if MaxQueryRange is given.
func RandomURL(qr QueryRange) (string, error) {
	return _FE_URL.Random(qr)
}

// generate source URL with given query.
func GenerateURL(q Query) (string, error) {
	return _FE_URL.Generate(q)
}

// helper function to consume error for url string.
// it will panic if error is not nil.
func MustURL(url string, err error) string {
	if err != nil {
		panic(err)
	}
	return url
}
