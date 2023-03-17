package model

import (
	"time"

	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
)

type Journey struct {
	ID            uint               `json:"id" db:"jou_id"`
	GroupID       uint               `json:"group_id" db:"jou_group_id"`
	Passengers    uint               `json:"passengers" db:"jou_passengers"`
	Status        enum.JourneyStatus `json:"status" db:"jou_status"`
	CreationDate  time.Time          `json:"creation_date" db:"jou_creation_date"`
	UpdateDate    time.Time          `json:"update_date" db:"jou_update_date"`
	DeleteDate    *time.Time         `json:"delete_date" db:"jou_delete_date"`
	CarAssignedID *uint              `json:"car_assigned_id" db:"jou_car_assigned"`
	CarAssigned   *Car               `json:"car" db:"car"`
}
