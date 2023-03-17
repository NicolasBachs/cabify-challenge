package entity

import (
	"time"

	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
)

type Journey struct {
	ID            uint               `json:"id"`
	GroupID       uint               `json:"group_id"`
	Passengers    uint               `json:"passengers"`
	Status        enum.JourneyStatus `json:"status"`
	CreationDate  time.Time          `json:"creation_date"`
	UpdateDate    time.Time          `json:"update_date"`
	DeleteDate    *time.Time         `json:"delete_date"`
	CarAssignedID *uint              `json:"car_assigned_id"`
	CarAssigned   *Car               `json:"car"`
}
