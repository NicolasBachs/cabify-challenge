package repositoryImpl

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/test/mocks"
)

func TestNewJourneyRepository(t *testing.T) {
	// Given
	dbx, _ := mocks.MockDatabase(t)
	defer dbx.Close()

	// Assert
	assert.NotNil(t, NewJourneyRepository(dbx))
}

func Test_JourneyRepository_GetJourneyByID_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()
	timeNow := time.Now()

	journeyID := uint(1)
	r := NewJourneyRepository(dbx)

	expectedResult := &entity.Journey{
		ID:            1,
		GroupID:       123,
		Passengers:    4,
		Status:        enum.JourneyStatusPending,
		CreationDate:  timeNow,
		UpdateDate:    timeNow,
		DeleteDate:    nil,
		CarAssignedID: nil,
	}

	rowsMock := sqlmock.NewRows(
		[]string{"jou_id", "jou_group_id", "jou_passengers", "jou_status", "jou_creation_date", "jou_update_date", "jou_delete_date", "jou_car_assigned"},
	).AddRow(1, 123, 4, enum.JourneyStatusPending, timeNow, timeNow, nil, nil)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(journeyID).WillReturnRows(rowsMock)
	actualResult, err := r.GetJourneyByID(journeyID)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_JourneyRepository_CreateJourney_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	journey := &entity.Journey{
		GroupID:       123,
		Passengers:    4,
		Status:        enum.JourneyStatusPending,
		CarAssignedID: nil,
	}
	r := NewJourneyRepository(dbx)

	expectedResult := uint(1)

	rowsMock := sqlmock.NewRows([]string{"jou_id"}).AddRow(1)

	// When
	mock.ExpectQuery("^INSERT (.+)").WithArgs(
		journey.GroupID, journey.Passengers, journey.Status, journey.CarAssignedID,
	).WillReturnRows(rowsMock)

	actualResult, err := r.CreateJourney(journey)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_JourneyRepository_GetJourneysByStatus_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()
	timeNow := time.Now()

	journeyStatus := enum.JourneyStatusPending
	r := NewJourneyRepository(dbx)

	expectedResult := []*entity.Journey{
		{
			ID:            1,
			GroupID:       101,
			Passengers:    3,
			Status:        enum.JourneyStatusPending,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: nil,
		},
		{
			ID:            2,
			GroupID:       102,
			Passengers:    2,
			Status:        enum.JourneyStatusPending,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: nil,
		},
		{
			ID:            3,
			GroupID:       103,
			Passengers:    1,
			Status:        enum.JourneyStatusPending,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: nil,
		},
	}

	rowsMock := sqlmock.NewRows(
		[]string{"jou_id", "jou_group_id", "jou_passengers", "jou_status", "jou_creation_date", "jou_update_date", "jou_delete_date", "jou_car_assigned"},
	).AddRow(
		1, 101, 3, enum.JourneyStatusPending, timeNow, timeNow, nil, nil,
	).AddRow(
		2, 102, 2, enum.JourneyStatusPending, timeNow, timeNow, nil, nil,
	).AddRow(
		3, 103, 1, enum.JourneyStatusPending, timeNow, timeNow, nil, nil,
	)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(journeyStatus).WillReturnRows(rowsMock)
	actualResult, err := r.GetJourneysByStatus(journeyStatus)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_JourneyRepository_CleanJourneys_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	r := NewJourneyRepository(dbx)

	// When
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE (.+)").WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	err := r.CleanJourneys()

	// Then
	assert.Nil(t, err)
}

func Test_JourneyRepository_CleanJourneys_RollbackOnFailLogicalDelete(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	r := NewJourneyRepository(dbx)

	// When
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE (.+)").WillReturnError(errors.New("some error"))
	mock.ExpectRollback()

	err := r.CleanJourneys()

	// Then
	assert.NotNil(t, err)
}

func Test_JourneyRepository_UpdateJourneyStatus(t *testing.T) {

	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	journeyID := uint(1)
	newStatus := enum.JourneyStatusFinished
	r := NewJourneyRepository(dbx)

	// When
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE (.+)").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := r.UpdateJourneyStatus(journeyID, newStatus)

	// Then
	assert.Nil(t, err)
}

func Test_JourneyRepository_GetPendingJourneysWherePassengersLessOrEqualTo_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()
	timeNow := time.Now()

	maxPassengers := uint(3)
	r := NewJourneyRepository(dbx)

	expectedResult := []*entity.Journey{
		{
			ID:            1,
			GroupID:       101,
			Passengers:    3,
			Status:        enum.JourneyStatusPending,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: nil,
		},
		{
			ID:            2,
			GroupID:       102,
			Passengers:    2,
			Status:        enum.JourneyStatusPending,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: nil,
		},
		{
			ID:            3,
			GroupID:       103,
			Passengers:    1,
			Status:        enum.JourneyStatusPending,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: nil,
		},
	}

	rowsMock := sqlmock.NewRows(
		[]string{"jou_id", "jou_group_id", "jou_passengers", "jou_status", "jou_creation_date", "jou_update_date", "jou_delete_date", "jou_car_assigned"},
	).AddRow(
		1, 101, 3, enum.JourneyStatusPending, timeNow, timeNow, nil, nil,
	).AddRow(
		2, 102, 2, enum.JourneyStatusPending, timeNow, timeNow, nil, nil,
	).AddRow(
		3, 103, 1, enum.JourneyStatusPending, timeNow, timeNow, nil, nil,
	)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(enum.JourneyStatusPending, maxPassengers).WillReturnRows(rowsMock)
	actualResult, err := r.GetPendingJourneysWherePassengersLessOrEqualTo(maxPassengers)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_JourneyRepository_GetPendingJourneysWherePassengersLessOrEqualTo_EmptyResult(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	maxPassengers := uint(3)
	r := NewJourneyRepository(dbx)

	var expectedResult []*entity.Journey

	rowsMock := sqlmock.NewRows(
		[]string{"jou_id", "jou_group_id", "jou_passengers", "jou_status", "jou_creation_date", "jou_update_date", "jou_delete_date", "jou_car_assigned"},
	)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(enum.JourneyStatusPending, maxPassengers).WillReturnRows(rowsMock)
	actualResult, err := r.GetPendingJourneysWherePassengersLessOrEqualTo(maxPassengers)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_JourneyRepository_AssignCarToJourney_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	journeyID := uint(1)
	carID := uint(2)
	r := NewJourneyRepository(dbx)

	// When
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE (.+)").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := r.AssignCarToJourney(journeyID, carID)

	// Then
	assert.Nil(t, err)
}

func Test_JourneyRepository_GetJourneysByCarID_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()
	timeNow := time.Now()

	carID := uint(3)
	r := NewJourneyRepository(dbx)

	expectedResult := []*entity.Journey{
		{
			ID:            1,
			GroupID:       101,
			Passengers:    1,
			Status:        enum.JourneyStatusAssigned,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: &carID,
		},
		{
			ID:            2,
			GroupID:       102,
			Passengers:    2,
			Status:        enum.JourneyStatusAssigned,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: &carID,
		},
		{
			ID:            3,
			GroupID:       103,
			Passengers:    3,
			Status:        enum.JourneyStatusAssigned,
			CreationDate:  timeNow,
			UpdateDate:    timeNow,
			DeleteDate:    nil,
			CarAssignedID: &carID,
		},
	}

	rowsMock := sqlmock.NewRows(
		[]string{"jou_id", "jou_group_id", "jou_passengers", "jou_status", "jou_creation_date", "jou_update_date", "jou_delete_date", "jou_car_assigned"},
	).AddRow(
		1, 101, 1, enum.JourneyStatusAssigned, timeNow, timeNow, nil, carID,
	).AddRow(
		2, 102, 2, enum.JourneyStatusAssigned, timeNow, timeNow, nil, carID,
	).AddRow(
		3, 103, 3, enum.JourneyStatusAssigned, timeNow, timeNow, nil, carID,
	)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(carID).WillReturnRows(rowsMock)
	actualResult, err := r.GetJourneysByCarID(carID)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_JourneyRepository_GetJourneysByCarID_EmptyResult(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	carID := uint(3)
	r := NewJourneyRepository(dbx)

	var expectedResult []*entity.Journey

	rowsMock := sqlmock.NewRows(
		[]string{"jou_id", "jou_group_id", "jou_passengers", "jou_status", "jou_creation_date", "jou_update_date", "jou_delete_date", "jou_car_assigned"},
	)
	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(carID).WillReturnRows(rowsMock)
	actualResult, err := r.GetJourneysByCarID(carID)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}
