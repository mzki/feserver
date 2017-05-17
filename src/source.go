package src

import "fmt"

// TODO: IntervalTime is contained in Source?
// it makes easy to construct new Getter.

// Default Source definitions. You can use these to
// construct new Getter by:
//	fe_getter := src.NewGetter(src.FE, src.LeastIntervalTime)
var (
	// source loacation for F.E. examination.
	FE = Source{
		URL: `http://www.fe-siken.com/kakomon/{{.Year}}_{{.Season}}/q{{.No}}.html`,
		QueryRange: QueryRange{
			MaxYear: 29, MinYear: 13,
			MaxNo: 80, MinNo: 1,
			Season: SeasonAll,
		},
	}

	// source location for A.P. examination.
	AP = Source{
		URL: `http://www.ap-siken.com/kakomon/{{.Year}}_{{.Season}}/q{{.No}}.html`,
		QueryRange: QueryRange{
			MaxYear: 29, MinYear: 13,
			MaxNo: 80, MinNo: 1,
			Season: SeasonAll,
		},
	}
)

// Source is the source definition for getting the questions
// from the extenal server.
type Source struct {
	URL string // URL template for source server.

	QueryRange // Acceptable range for query.
}

// check whether itself has correct values?
// return nil if it is valid.
func (src Source) ValidatesSelf() error {
	if src.MaxYear < src.MinYear {
		return fmt.Errorf("Source: MaxYear must be larger then MinYear but Max: %d, Min: %d", src.MaxYear, src.MinYear)
	}
	if src.MaxNo < src.MinNo {
		return fmt.Errorf("Source: MaxNo must be larger then MinNo but Max: %d, Min: %d", src.MaxNo, src.MinNo)
	}
	switch src.Season {
	case SeasonSpring, SeasonAutumn, SeasonAll:
		return nil
	default:
		return fmt.Errorf("Source: Season must be either %s, %s or %s", SeasonSpring, SeasonAutumn, SeasonAll)
	}
}

// check whether given query has correct value range
// in the source? nil error means query is valid.
func (src Source) Validates(q Query) error {
	if y := q.Year; y < src.MinYear || y > src.MaxYear {
		return fmt.Errorf("Query: year must be in [%d:%d], but %d", src.MinYear, src.MaxYear, y)
	}
	if n := q.No; n < src.MinNo || n > src.MaxNo {
		return fmt.Errorf("Query: Question No. must be in [%d:%d], but %d", src.MinNo, src.MaxNo, n)
	}

	qSeason := q.Season
	switch {
	case qSeason != SeasonSpring && qSeason != SeasonAutumn:
		return fmt.Errorf("Query: Season must be either %s or %s, but %s", SeasonSpring, SeasonAutumn, qSeason)
	case src.Season == SeasonSpring && qSeason != SeasonSpring:
		return fmt.Errorf("Query: Season must be %s but %s", SeasonSpring, qSeason)
	case src.Season == SeasonAutumn && qSeason != SeasonAutumn:
		return fmt.Errorf("Query: Season must be %s but %s", SeasonAutumn, qSeason)
	}
	return nil
}

// check whether QueryRange is in correct range in the Source?
// return nil if query is valid.
func (src Source) ValidatesRange(qr QueryRange) error {
	{
		// qr.Season is validated later. dummy Season is used insteadly.
		dummySeason := SeasonSpring
		if src.Season == SeasonAutumn {
			dummySeason = SeasonAutumn
		}
		if err := src.Validates(Query{qr.MaxYear, dummySeason, qr.MaxNo}); err != nil {
			return err
		}
		if err := src.Validates(Query{qr.MinYear, dummySeason, qr.MinNo}); err != nil {
			return err
		}
	}

	// check the relation for min and max.
	if qr.MaxYear < qr.MinYear {
		return fmt.Errorf("QueryRange: MaxYear must be larger then MinYear but Max: %d, Min: %d", qr.MaxYear, qr.MinYear)
	}
	if qr.MaxNo < qr.MinNo {
		return fmt.Errorf("QueryRange: MaxNo must be larger then MinNo but Max: %d, Min: %d", qr.MaxNo, qr.MinNo)
	}
	// check season
	qSeason := qr.Season
	switch {
	case qSeason != SeasonSpring && qSeason != SeasonAutumn && qSeason != SeasonAll:
		return fmt.Errorf("QueryRange: Season must be either %s, %s or %s, but %s",
			SeasonSpring, SeasonAutumn, SeasonAll, qSeason)
	case src.Season == SeasonSpring && qSeason != SeasonSpring:
		return fmt.Errorf("Query: Season must be %s but %s", SeasonSpring, qSeason)
	case src.Season == SeasonAutumn && qSeason != SeasonAutumn:
		return fmt.Errorf("Query: Season must be %s but %s", SeasonAutumn, qSeason)
	}
	return nil
}
