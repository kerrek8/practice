package models

type Listing struct {
	ID          int64
	UserID      int64
	Name        string
	Price       float64
	City        string
	Description string
	Status      string
}
