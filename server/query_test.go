package server

import (
	"fmt"
	"net/url"
	"testing"
)

func TestParseGetRandomQuery(t *testing.T) {
	const (
		MaxYear = 20
		MinYear = 20
		MaxNo   = 10
		MinNo   = 10
		Season  = "all"
	)
	url, err := url.Parse(fmt.Sprintf("localhost?max_year=%d&min_year=%d&season=%s&max_no=%d&min_no=%d",
		MaxYear, MinYear, Season, MaxNo, MinNo))
	if err != nil {
		t.Fatal(err)
	}

	qr, err := parseGetRandomQuery(url.Query(), DefaultSource)
	if err != nil {
		t.Fatal(err)
	}

	assertEqualInt(t, qr.MaxYear, MaxYear, "")
	assertEqualInt(t, qr.MinYear, MinYear, "")
	assertEqualInt(t, qr.MinNo, MinNo, "")
	assertEqualInt(t, qr.MaxNo, MaxNo, "")
}

func assertEqualInt(t *testing.T, got, expect int, mes string) {
	if got != expect {
		t.Errorf("must be equal but got: %d, expect: %d, "+mes, got, expect)
	}
}
