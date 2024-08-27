package model

type (
	Division struct {
		ID   string
		Name string       `json:"name" validate:"required,min=3,max=32"`
		Age  DivisionAges `json:"age" validate:"required"`
	}

	DivisionAges struct {
		Min int `json:"min"`
		Max int `json:"max"`
	}

	DivisionCreation struct {
		Divisions []Division `json:"divisions" validate:"required"`
	}
)
