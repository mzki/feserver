package src

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestGetRandom(t *testing.T) {
	t.Skip("Server Access is prevented")
	// because the server stands outsider, accessing to the server should not be frequent.
	const N = 10
	const WAIT = 1
	ctx := context.Background()
	for i := 0; i < N; i++ {
		time.Sleep(time.Duration(rand.Intn(WAIT)+1) * time.Second) // wait time for server access.
		res, err := GetRandom(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.URL) == 0 {
			t.Fatal("empty response returned")
		}

		var any_not_parsed = false
		if len(res.Question) == 0 {
			any_not_parsed = true
			t.Error("question is not parsed")
		}
		if len(res.Answer) == 0 {
			any_not_parsed = true
			t.Error("answer is not parsed")
		}
		if len(res.Selections) == 0 {
			any_not_parsed = true
			t.Error("selections are not parsed")
		}

		if any_not_parsed {
			t.Errorf("not parsed URL: %s", res.URL)
		}
	}
}

func TestGet(t *testing.T) {
	t.Skip("Server Access is prevented")
	res, err := Get(context.Background(), Query{
		Season: SeasonSpring,
		Year:   28,
		No:     2,
	})
	if err != nil {
		t.Fatal(err)
	}

	fpout, err := os.Create("got_out.html")
	defer fpout.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprintf(fpout, "Question:\n%v\nSelections:\n%v\nAnswer:\n%v\nExplanation:\n%v\n",
		res.Question, res.Selections, res.Answer, res.Explanation)
}

func TestRandomQuery(t *testing.T) {
	const randomN = 100
	for i := 0; i < randomN; i++ {
		q := randomQuery(&DefaultQueryRange)
		if s := q.Season; s != seasonRange[0] && s != seasonRange[1] {
			t.Errorf("invaid season, got: %s", s)
		}
		if y := q.Year; y < yearRange[0] || y > yearRange[1] {
			t.Errorf("invaid year, got: %v", y)
		}
		if n := q.No; n < noRange[0] || n > noRange[1] {
			t.Errorf("invaid number, got: %v", n)
		}
	}
}

func TestParseDoc(t *testing.T) {
	fp, err := os.Open("./test.html")
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()

	doc, err := goquery.NewDocumentFromReader(fp)
	if err != nil {
		t.Fatal(err)
	}

	res, err := parseDoc(doc)
	if err != nil {
		t.Fatal(err)
	}

	fpout, err := os.Create("out.html")
	defer fpout.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprintf(fpout, "Question:\n%v\nSelections:\n%v\nAnswer:\n%v\nExplanation:\n%v\n",
		res.Question, res.Selections, res.Answer, res.Explanation)
}

func TestParseDocHasImage(t *testing.T) {
	fp, err := os.Open("./y19_spring_q26.html")
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()

	doc, err := goquery.NewDocumentFromReader(fp)
	if err != nil {
		t.Fatal(err)
	}

	res, err := parseDoc(doc)
	if err != nil {
		t.Fatal(err)
	}

	if res.HasImage {
		t.Fatal("must not have some image, but HasImage is true")
	}
}
