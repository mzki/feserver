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
		res, err := GetRandom(ctx, MaxQueryRange)
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

	fpout, err := os.Create("test_get.txt")
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
		q := randomQuery(FE.QueryRange)
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
	doc, err := goqueryDocFile("./y28_spring_q2.html")
	if err != nil {
		t.Fatal(err)
	}

	res, err := parseDoc(doc)
	if err != nil {
		t.Fatal(err)
	}

	fpout, err := os.Create("test_parse_doc.txt")
	defer fpout.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprintf(fpout, "Question:\n%v\nSelections:\n%v\nAnswer:\n%v\nExplanation:\n%v\n",
		res.Question, res.Selections, res.Answer, res.Explanation)
}

func TestParseDocHasImage(t *testing.T) {
	doc, err := goqueryDocFile("./y19_spring_q26.html")
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

func TestParseDocExplainCh(t *testing.T) {
	for _, testcase := range []struct {
		in, out string
	}{
		{"./y28_spring_q2.html", "./y28_spring_q2.txt"},
		{"./y19_spring_q26.html", "./y19_spring_q26.txt"},
	} {
		doc, err := goqueryDocFile(testcase.in)
		if err != nil {
			t.Fatal(err)
		}

		res, err := parseDoc(doc)
		if err != nil {
			t.Fatal(err)
		}

		fpout, err := os.Create(testcase.out)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(fpout, "Question:\n%v\nSelections:\n%v\nAnswer:\n%v\nExplanation:\n%v\n",
			res.Question, res.Selections, res.Answer, res.Explanation)
		fpout.Close()
	}
}

func goqueryDocFile(file string) (*goquery.Document, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return goquery.NewDocumentFromReader(fp)
}

func BenchmarkParseDoc(b *testing.B) {
	doc, err := goqueryDocFile("./y28_spring_q2.html")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		clone := doc.Selection.Clone()
		b.StartTimer()
		_, err := parseDoc(doc)
		if err != nil {
			b.Fatal(err)
		}
		doc.Selection = clone
	}
}

func BenchmarkRandomURL(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := RandomURL(MaxQueryRange)
		if err != nil {
			b.Fatal(err)
		}
	}
}
