package model

type (
	Player struct {
		ID           string  `json:"id"`
		ParentID     string  `json:"parentID"`
		FirstName    string  `json:"firstName" validate:"required,min=2,max=32"`
		LastName     string  `json:"lastName" validate:"required,min=2,max=32"`
		DateOfBirth  string  `json:"dateOfBirth" validate:"required"`
		DivisionID   *string `json:"divisionID"`
		TeamID       *string `json:"teamID"`
		IsRegistered bool    `json:"isRegistered"`
		Hash         string  `json:"hash"`
	}

	PlayerIDs struct {
		Players []string `json:"players" validate:"required"`
	}

	PlayerList struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	Players struct {
		Players []Player `json:"players" validate:"required"`
	}
)
