package src

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

var targetURLTmpl = template.Must(template.New("targetURL").Parse(
	`http://www.fe-siken.com/kakomon/{{.Year}}_{{.Season}}/q{{.No}}.html`))

type query struct {
	Year   int
	Season string
	No     int
}

var (
	yearRange   = [...]int{15, 27}
	seasonRange = [...]string{"haru", "aki"}
	noRange     = [...]int{1, 80}
)

// generates random query.
func randomQuery() query {
	year := rand.Intn(yearRange[1]-yearRange[0]+1) + yearRange[0]
	no := rand.Intn(noRange[1]-noRange[0]+1) + noRange[0]
	season := seasonRange[rand.Intn(len(seasonRange))]
	return query{year, season, no}
}

// generate random target URL.
func randomURL() string {
	q := randomQuery()
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
	Answer     string

	// indicates Question, Selections or Answer contain some image.
	// the response can not be represented by plain text only.
	HasImage bool

	URL string // source URL
}

// GetRandom returns a response, which is randomly selected F.E question and its answer, from website.
// This process takes some time. You can cancel it by canceling context.
func GetRandom(ctx context.Context) (Response, error) {
	resCh := make(chan Response, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(resCh)
		defer close(errCh)

		url := randomURL()
		doc, err := goquery.NewDocument(url)
		if err != nil {
			errCh <- err
			return
		}
		res, err := parseDoc(doc)
		if err != nil {
			errCh <- err
			return
		}
		res.URL = doc.Url.String()
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
		Question:   q_doc.Text(),
		Selections: selections,
		Answer:     fmt.Sprintf("正解：%s, 解説：%s", ansch_doc.Text(), ansbg_doc.Text()),
		HasImage:   has_image,
	}, nil
}
