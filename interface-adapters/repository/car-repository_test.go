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

func TestNewCarRepository(t *testing.T) {
	// Given
	dbx, _ := mocks.MockDatabase(t)
	defer dbx.Close()

	// Assert
	assert.NotNil(t, NewCarRepository(dbx))
}

func Test_CarRepository_GetCarByID_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()
	timeNow := time.Now()

	carID := 1
	r := NewCarRepository(dbx)

	expectedResult := &entity.Car{
		ID:             1,
		MaxSeats:       4,
		AvailableSeats: 2,
		CreationDate:   timeNow,
		UpdateDate:     timeNow,
		DeleteDate:     nil,
	}

	rowsMock := sqlmock.NewRows(
		[]string{"car_id", "car_max_seats", "car_creation_date", "car_update_date", "car_delete_date", "car_passengers_aboard"},
	).AddRow(1, 4, timeNow, timeNow, nil, 2)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(enum.JourneyStatusAssigned, carID).WillReturnRows(rowsMock)
	actualResult, err := r.GetCarByID(uint(carID))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_CarRepository_GetCarByID_SuccessEmptyResult(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	r := NewCarRepository(dbx)

	carID := 1

	rowsMock := sqlmock.NewRows(
		[]string{"car_id", "car_max_seats", "car_creation_date", "car_update_date", "car_delete_date", "car_passengers_aboard"},
	)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(enum.JourneyStatusAssigned, carID).WillReturnRows(rowsMock)
	actualResult, err := r.GetCarByID(uint(carID))

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, actualResult)
}

func Test_CarRepository_CreateCar_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	r := NewCarRepository(dbx)

	carMock := &entity.Car{
		ID:       100,
		MaxSeats: 5,
	}

	expectedResult := uint(1)

	rowsMock := sqlmock.NewRows([]string{"car_id"}).AddRow(1)

	// When
	mock.ExpectQuery("^INSERT INTO (.+)").WithArgs(carMock.MaxSeats).WillReturnRows(rowsMock)
	actualResult, err := r.CreateCar(carMock)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_CarRepository_DeleteCarsAndInsert_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	r := &carRepository{
		db: dbx,
	}

	carsMock := []*entity.Car{
		{
			MaxSeats: 4,
		},
		{
			MaxSeats: 3,
		},
		{
			MaxSeats: 2,
		},
	}

	expectedResult := []uint{1, 2, 3}

	// When
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE (.+)").WillReturnResult(sqlmock.NewResult(0, 2))
	for index := range carsMock {
		rows := sqlmock.NewRows([]string{"car_id"})
		rows.AddRow(index + 1)
		mock.ExpectQuery("^INSERT (.+)").WillReturnRows(rows)
	}
	mock.ExpectCommit()

	actualResult, err := r.DeleteCarsAndInsert(carsMock)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_CarRepository_DeleteCarsAndInsert_RollbackOnFailLogicalDelete(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	r := &carRepository{
		db: dbx,
	}

	carsMock := []*entity.Car{
		{
			MaxSeats: 4,
		},
		{
			MaxSeats: 3,
		},
		{
			MaxSeats: 2,
		},
	}

	// When
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE (.+)").WillReturnError(errors.New("Some error"))
	mock.ExpectRollback()

	actualResult, err := r.DeleteCarsAndInsert(carsMock)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, actualResult)
}

func Test_CarRepository_DeleteCarsAndInsert_RollbackOnFailInsert(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()

	r := &carRepository{
		db: dbx,
	}

	carsMock := []*entity.Car{
		{
			MaxSeats: 4,
		},
		{
			MaxSeats: 3,
		},
		{
			MaxSeats: 2,
		},
	}

	// When
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE (.+)").WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec("^INSERT (.+)").WillReturnError(errors.New("Some error"))
	mock.ExpectRollback()

	actualResult, err := r.DeleteCarsAndInsert(carsMock)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, actualResult)
}

func Test_CarRepository_GetCarsWithAvailableSeatsGreaterOrEqualTo_Success(t *testing.T) {
	// Given
	dbx, mock := mocks.MockDatabase(t)
	defer dbx.Close()
	timeNow := time.Now()

	r := &carRepository{
		db: dbx,
	}

	passengers := uint(2)

	expectedResult := []*entity.Car{
		{
			ID:             1,
			MaxSeats:       4,
			CreationDate:   timeNow,
			UpdateDate:     timeNow,
			DeleteDate:     nil,
			AvailableSeats: 2,
		},
		{
			ID:             1,
			MaxSeats:       4,
			CreationDate:   timeNow,
			UpdateDate:     timeNow,
			DeleteDate:     nil,
			AvailableSeats: 3,
		},
		{
			ID:             1,
			MaxSeats:       4,
			CreationDate:   timeNow,
			UpdateDate:     timeNow,
			DeleteDate:     nil,
			AvailableSeats: 4,
		},
	}

	rowsMock := sqlmock.NewRows(
		[]string{"car_id", "car_max_seats", "car_creation_date", "car_update_date", "car_delete_date", "car_passengers_aboard"},
	).AddRow(
		1, 4, timeNow, timeNow, nil, 2,
	).AddRow(
		1, 4, timeNow, timeNow, nil, 1,
	).AddRow(
		1, 4, timeNow, timeNow, nil, 0,
	)

	// When
	mock.ExpectQuery("^SELECT (.+)").WithArgs(enum.JourneyStatusAssigned, passengers).WillReturnRows(rowsMock)

	actualResult, err := r.GetCarsWithAvailableSeatsGreaterOrEqualTo(passengers)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}
