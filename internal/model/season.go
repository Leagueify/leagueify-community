package model

type (
	Season struct {
		ID           string
		Name         string      `json:"name" validate:"required"`
		Season       SeasonDates `json:"season" validate:"required"`
		Registration SeasonDates `json:"registration" validate:"required"`
	}

	SeasonDates struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	SeasonList struct {
		ID   string
		Name string
	}

	SeasonUpdate struct {
		ID           string
		Name         string
		Season       SeasonDates
		Registration SeasonDates
	}
)
