package resources

import "fmt"

type MutexResource string

var (
	ANY_CAR     MutexResource = "ANY_CAR"
	ANY_JOURNEY MutexResource = "ANY_JOURNEY"
)

func CarWithID(carID uint) MutexResource {
	return MutexResource(fmt.Sprintf("CAR_%d", carID))
}

func GroupWithID(groupID uint) MutexResource {
	return MutexResource(fmt.Sprintf("GROUP_%d", groupID))
}

func JourneyWithID(journeyID uint) MutexResource {
	return MutexResource(fmt.Sprintf("JOURNEY_%d", journeyID))
}

func (r MutexResource) ToString() string {
	return string(r)
}
