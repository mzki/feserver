//go:generate sh ./gen_test.sh

package src

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Response is a result of parsing a web page
// that have the F.E. question and its answer.
type Response struct {
	Question   string   `json:"question"`
	Selections []string `json:"selections"`

	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`

	// indicates Question, Selections or Answer contain some image.
	// the response can not be represented by plain text only.
	HasImage bool `json:"hasImage"`

	URL string `json:"url"` // source URL
}

var defaultGetter = NewGetter(LeastIntervalTime)

// the minimum time for request interval.
const LeastIntervalTime = 5 * time.Second

// interval time varies plus or minus VariationCoef second.
const VariationCoef = 2

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
		panic("intervalTime must be >= " + LeastIntervalTime.String())
	}
	return &Getter{intervalTime: intervalTime, lastRequest: time.Time{}}
}

func (g *Getter) wait() {
	// all of the exported method check this.
	if g.intervalTime < LeastIntervalTime {
		panic("intervalTime must be >= " + LeastIntervalTime.String())
	}

	randMutex.Lock()
	coef := random.Intn(2*VariationCoef+1) - VariationCoef
	randMutex.Unlock()

	noisedInterval := g.intervalTime + time.Duration(coef)*time.Second
	if wait := noisedInterval - time.Since(g.lastRequest); wait > 0 {
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
// use default random range if RandomRange is nil.
func (g *Getter) GetRandom(ctx context.Context, qr *QueryRange) (Response, error) {
	g.wait()
	return getResponse(ctx, RandomURL(qr))
}

// Get() returns a response, which contains F.E question and its answer selected by Query, from website.
// This process takes some time. You can cancel it by canceling context.
func Get(ctx context.Context, q Query) (Response, error) {
	return defaultGetter.Get(ctx, q)
}

// GetRandom() returns a response, which is randomly selected F.E question and its answer, from website.
// This process takes some time. You can cancel it by canceling context.
// use default random range if RandomRange is nil.
func GetRandom(ctx context.Context, qr *QueryRange) (Response, error) {
	return defaultGetter.GetRandom(ctx, qr)
}

func getResponse(ctx context.Context, url string) (Response, error) {
	resCh := make(chan Response, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(resCh)
		defer close(errCh)
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
// Because target url has a content encoded by ShiftJIS,
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
	if q_doc.Find("img").Length() > 0 || sel_doc.Find("img").Length() > 0 || ans_doc.Find("img").Length() > 0 {
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
