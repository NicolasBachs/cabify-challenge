package model

import "time"

type Car struct {
	ID               uint       `json:"id" db:"car_id"`
	MaxSeats         uint       `json:"maxSeats" db:"car_max_seats"`
	CreationDate     time.Time  `json:"creation_date" db:"car_creation_date"`
	UpdateDate       time.Time  `json:"update_date" db:"car_update_date"`
	DeleteDate       *time.Time `json:"delete_date" db:"car_delete_date"`
	PassengersAboard uint       `json:"passengers_aboard" db:"car_passengers_aboard"`
}
