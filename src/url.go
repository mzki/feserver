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

// Query is a query for source URL.
// Its fields specifies which question is searched for.
type Query struct {
	Year   int
	Season string
	No     int
}

const (
	// These represents query range.

	// Examination Year range
	MinYear = 15
	MaxYear = 28

	// Examination Season selections
	SeasonSpring = "haru"
	SeasonAutumn = "aki"

	// Question No.
	MinNO = 1
	MaxNO = 80
)

var (
	yearRange   = [...]int{MinYear, MaxYear}
	seasonRange = [...]string{SeasonSpring, SeasonAutumn}
	noRange     = [...]int{MinNO, MaxNO}
)

// check whether query has correct value range?
// nil error means query is valid.
func (q Query) Validates() error {
	if y := q.Year; y < MinYear || y > MaxYear {
		return fmt.Errorf("year must be in [%d:%d], but %d", MinYear, MaxYear, y)
	}
	if s := q.Season; s != SeasonSpring && s != SeasonAutumn {
		return fmt.Errorf("Season must be either %s or %s, but %s", SeasonSpring, SeasonAutumn, s)
	}
	if n := q.No; n < MinNO || n > MaxNO {
		return fmt.Errorf("Question No. must be in [%d:%d], but %d", MinNO, MaxNO, n)
	}
	return nil
}

var randMutex = new(sync.Mutex)

// package global random state. under mutex.
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// generates random query.
func randomQuery() Query {
	randMutex.Lock()
	year := random.Intn(MaxYear-MinYear+1) + MinYear
	no := random.Intn(MaxNO-MinNO+1) + MinNO
	season := seasonRange[random.Intn(len(seasonRange))]
	randMutex.Unlock()
	return Query{year, season, no}
}

// generate target URL with randomized query.
func RandomURL() string {
	q := randomQuery()
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
