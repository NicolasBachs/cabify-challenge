package repositoryImpl

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/entity"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/enum"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/repository"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/mapper"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/model"
)

type carRepository struct {
	db *sqlx.DB
}

func NewCarRepository(
	db *sqlx.DB,
) repository.CarRepository {

	return &carRepository{
		db: db,
	}
}

func (r *carRepository) GetCarByID(carID uint) (*entity.Car, error) {
	var carModel model.Car

	err := r.db.Get(
		&carModel,
		`SELECT 
			c.car_id,
			c.car_max_seats,
			c.car_creation_date,
			c.car_update_date,
			c.car_delete_date,
			SUM(CASE WHEN j.jou_status <> $1 OR j.jou_passengers IS NULL THEN 0 ELSE j.jou_passengers END) as car_passengers_aboard
		FROM cars c 
		LEFT JOIN journeys j ON j.jou_car_assigned = c.car_id 
		WHERE 
			c.car_delete_date IS NULL
			AND j.jou_delete_date IS NULL
			AND c.car_id = $2
		GROUP BY c.car_id, c.car_max_seats, c.car_creation_date, c.car_update_date, c.car_delete_date`,
		enum.JourneyStatusAssigned, carID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	carEntity := mapper.CarModelToEntity(&carModel)
	carEntity.AvailableSeats = carEntity.MaxSeats - carModel.PassengersAboard

	return carEntity, nil
}

func (r *carRepository) CreateCar(car *entity.Car) (uint, error) {
	var id uint

	err := r.db.QueryRowx(`
        INSERT INTO cars (car_max_seats, car_creation_date, car_update_date)
        VALUES ($1, NOW(), NOW())
        RETURNING car_id
    `, car.MaxSeats).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Reset cars, clean existing cars and create multiple cars in the same transaction
func (r *carRepository) DeleteCarsAndInsert(cars []*entity.Car) ([]uint, error) {
	var carIDs []uint
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		DELETE FROM cars 
		WHERE 1 = 1
	`)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, car := range cars {
		var carID uint

		err := tx.QueryRowx(`
			INSERT INTO cars (car_id, car_max_seats, car_creation_date, car_update_date)
			VALUES ($1, $2, NOW(), NOW())
			RETURNING car_id
		`, car.ID, car.MaxSeats).Scan(&carID)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		carIDs = append(carIDs, carID)
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return carIDs, nil
}

func (r *carRepository) GetCarsWithAvailableSeatsGreaterOrEqualTo(passengers uint) ([]*entity.Car, error) {
	var cars []*model.Car

	err := r.db.Select(
		&cars,
		`SELECT 
			c.car_id,
			c.car_max_seats,
			c.car_creation_date,
			c.car_update_date,
			c.car_delete_date,
			SUM(CASE WHEN j.jou_status <> $1 OR j.jou_passengers IS NULL THEN 0 ELSE j.jou_passengers END) as car_passengers_aboard
		FROM cars c 
		LEFT JOIN journeys j ON j.jou_car_assigned = c.car_id 
		WHERE 
			c.car_delete_date IS NULL
			AND j.jou_delete_date IS NULL
		GROUP BY c.car_id, c.car_max_seats, c.car_creation_date, c.car_update_date, c.car_delete_date
		HAVING (c.car_max_seats - SUM(CASE WHEN j.jou_status <> $1 OR j.jou_passengers IS NULL THEN 0 ELSE j.jou_passengers END)) >= $2
		ORDER BY (c.car_max_seats - SUM(CASE WHEN j.jou_status <> $1 OR j.jou_passengers IS NULL THEN 0 ELSE j.jou_passengers END)) ASC`,
		enum.JourneyStatusAssigned, passengers,
	)

	var result []*entity.Car

	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}
		return nil, err
	}

	for _, carModel := range cars {
		carEntity := mapper.CarModelToEntity(carModel)
		carEntity.AvailableSeats = carEntity.MaxSeats - carModel.PassengersAboard
		result = append(result, carEntity)
	}

	return result, nil
}
