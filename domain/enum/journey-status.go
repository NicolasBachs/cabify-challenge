package enum

type JourneyStatus string

const (
	JourneyStatusPending   JourneyStatus = "PENDING"
	JourneyStatusAssigned  JourneyStatus = "ASSIGNED"
	JourneyStatusFinished  JourneyStatus = "FINISHED"
	JourneyStatusCancelled JourneyStatus = "CANCELLED"
)
