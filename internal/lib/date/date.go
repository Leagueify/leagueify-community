package date

import (
	"time"

	"github.com/leagueify/leagueify/internal/lib/error"
)

// comparisonDate is current date if not provided
func DifferenceInYears(providedDate string, comparisionDate *string) (int, error) {
	pDate, err := time.Parse(time.DateOnly, providedDate)
	if err != nil {
		return 0, err
	}

	if comparisionDate == nil {
		comparisionDate = new(string)
		*comparisionDate = time.Now().Format(time.DateOnly)
	}

	cDate, err := time.Parse(time.DateOnly, *comparisionDate)
	if err != nil {
		return 0, err
	}

	yearsDiff := cDate.Year() - pDate.Year()
	if cDate.Month() < pDate.Month() ||
		(cDate.Month() == pDate.Month() && cDate.Day() < pDate.Day()) {
		yearsDiff--
	}
	if yearsDiff < 0 {
		return 0, &errors.LeagueifyError{
			Message: "negative difference in years",
		}
	}
	return yearsDiff, nil
}

// comparisonDate is current date if not provided
func MeetsYearRequirement(years int, providedDate string, comparisionDate *string) bool {
	difference, err := DifferenceInYears(providedDate, comparisionDate)
	if err != nil {
		return false
	}

	if difference >= years {
		return true
	}

	return false
}
