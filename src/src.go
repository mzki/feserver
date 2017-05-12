package src

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

// TODO: has image logic is incorrect.

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

// generates random query.
func randomQuery() Query {
	year := rand.Intn(MaxYear-MinYear+1) + MinYear
	no := rand.Intn(MaxNO-MinNO+1) + MinNO
	season := seasonRange[rand.Intn(len(seasonRange))]
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

// Response is a result of parsing a web page
// that have the F.E. question and its answer.
type Response struct {
	Question   string
	Selections []string

	Answer      string
	Explanation string

	// indicates Question, Selections or Answer contain some image.
	// the response can not be represented by plain text only.
	HasImage bool

	URL string // source URL
}

var defaultGetter = NewGetter(LeastIntervalTime)

// the minimum time for request interval.
const LeastIntervalTime = 5 * time.Second

// Getter is a interface for F.E. question and answer from webpage.
// serial requests are splited by some interval time so that
// the number of accessing the outer server is reduced.
type Getter struct {
	intervalTime time.Duration
	lastRequest  time.Time
}

// return new Getter with intervalTime for server request.
func NewGetter(intervalTime time.Duration) *Getter {
	if intervalTime < LeastIntervalTime {
		panic("intervalTime must be > " + LeastIntervalTime.String())
	}
	return &Getter{intervalTime: intervalTime, lastRequest: time.Time{}}
}

func (g *Getter) wait() {
	if wait := g.intervalTime - time.Since(g.lastRequest); wait > 0 {
		time.Sleep(wait)
	}
	g.lastRequest = time.Now()
}

// Get() returns a response, which contains F.E question and its answer selected by Query, from website.
// This process takes some time. You can cancel it by canceling context.
//
// The interval wait time is inserted between serial calling of this method.
func (g *Getter) Get(ctx context.Context, q Query) (Response, error) {
	g.wait()
	return getResponse(ctx, GenerateURL(q))
}

// Get() returns a response, which contains F.E question and its answer selected by Query, from website.
// This process takes some time. You can cancel it by canceling context.
//
// The interval wait time is inserted between serial calling of this method.
func (g *Getter) GetRandom(ctx context.Context) (Response, error) {
	g.wait()
	return getResponse(ctx, RandomURL())
}

// Get() returns a response, which contains F.E question and its answer selected by Query, from website.
// This process takes some time. You can cancel it by canceling context.
func Get(ctx context.Context, q Query) (Response, error) {
	return defaultGetter.Get(ctx, q)
}

// GetRandom() returns a response, which is randomly selected F.E question and its answer, from website.
// This process takes some time. You can cancel it by canceling context.
func GetRandom(ctx context.Context) (Response, error) {
	return defaultGetter.GetRandom(ctx)
}

func getResponse(ctx context.Context, url string) (Response, error) {
	resCh := make(chan Response, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(resCh)
		defer close(errCh)
		// doc, err := goquery.NewDocument(url)
		doc, err := newDocument(url)
		if err != nil {
			errCh <- err
			return
		}
		res, err := parseDoc(doc)
		if err != nil {
			errCh <- err
			return
		}
		res.URL = url
		resCh <- res
	}()

	select {
	case <-ctx.Done():
		return Response{}, ctx.Err()
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return Response{}, err
	}
}

// newDocument() returns goquery.Document with UTF8 form.
//
// Because target url encoded by ShiftJIS,
// conversion from ShiftJIS to UTF8 is required before parsing goquery.Document.
// newDocument() performs that.
func newDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	r := transform.NewReader(bytes.NewReader(html), japanese.ShiftJIS.NewDecoder())
	return goquery.NewDocumentFromReader(r)
}

func parseDoc(doc *goquery.Document) (Response, error) {
	// parse section for a question.
	q_doc := doc.Find("div.main.kako > h3.qno").Next()
	if q_doc == nil {
		panic("nil Document")
	}

	// parse selections for a answer.
	sel_doc := doc.Find("div.main.kako > div.ansbg > ul.selectList.cf")
	if sel_doc == nil {
		panic("nil Document")
	}
	selections := sel_doc.Children().Filter("li").Map(func(_ int, s *goquery.Selection) string {
		sel_text := s.Find("div").Text()
		sel_ch := s.Find("a.selectBtn > button").Text()
		return sel_ch + ": " + sel_text
	})

	// parse section for a answer.
	ans_doc := doc.Find("div.main.kako > div.answerBox")
	if ans_doc == nil {
		panic("nil Document")
	}
	ansch_doc := ans_doc.Find("span#answerChar")
	if ansch_doc == nil {
		panic("nil Document")
	}
	ansbg_doc := ans_doc.Next().Next() // div.ansbg
	if ansbg_doc == nil {
		panic("nil Document")
	}

	// check whether the question has some image?
	var has_image = false
	if q_doc.Find("img") != nil || sel_doc.Find("img") != nil || ans_doc.Find("img") != nil {
		has_image = true
	}

	return Response{
		Question:    q_doc.Text(),
		Selections:  selections,
		Answer:      ansch_doc.Text(),
		Explanation: ansbg_doc.Text(),
		HasImage:    has_image,
	}, nil
}

// ParseHTML is helper funtion which parses html text and
// converts to f.e. question Response.
func ParseHTML(html string) (Response, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return Response{}, err
	}

	return parseDoc(doc)
}
