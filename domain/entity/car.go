package entity

import "time"

type Car struct {
	ID             uint       `json:"id"`
	MaxSeats       uint       `json:"seats"`
	AvailableSeats uint       `json:"availableSeats"`
	CreationDate   time.Time  `json:"creation_date"`
	UpdateDate     time.Time  `json:"update_date"`
	DeleteDate     *time.Time `json:"delete_date"`
	Journeys       []*Journey `json:"journeys"`
}
