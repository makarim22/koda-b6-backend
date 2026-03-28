package models

type Temperature struct {
	ID    int    `json:"id" db:"id"`
	Label string `json:"label" db:"label"`
}
