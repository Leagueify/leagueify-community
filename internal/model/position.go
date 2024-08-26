package model

type (
	Position struct {
		ID   string
		Name string
	}
	PositionCreation struct {
		Positions []string `json:"positions" validate:"required"`
	}
)
