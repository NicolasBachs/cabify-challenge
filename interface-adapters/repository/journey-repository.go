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

type journeyRepository struct {
	db *sqlx.DB
}

func NewJourneyRepository(db *sqlx.DB) repository.JourneyRepository {
	return &journeyRepository{
		db: db,
	}
}

func (r *journeyRepository) GetJourneyByID(journeyID uint) (*entity.Journey, error) {
	var journey model.Journey
	err := r.db.Get(
		&journey,
		"SELECT * FROM journeys WHERE jou_id = $1 AND jou_delete_date IS NULL",
		journeyID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return mapper.JourneyModelToEntity(&journey), nil
}

func (r *journeyRepository) GetLastJourneyByGroupID(groupID uint) (*entity.Journey, error) {
	var journey model.Journey
	err := r.db.Get(
		&journey,
		`
		SELECT 
			DISTINCT ON (jou_id) jou_id, 
			jou_group_id, 
			jou_passengers, 
			jou_status, 
			jou_car_assigned, 
			jou_creation_date, 
			jou_update_date
		FROM journeys
		WHERE
			jou_group_id = $1 
			AND jou_delete_date IS NULL
		ORDER BY jou_id, jou_creation_date DESC;
    	`, groupID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return mapper.JourneyModelToEntity(&journey), nil
}

func (r *journeyRepository) CreateJourney(journey *entity.Journey) (uint, error) {
	var id uint

	err := r.db.QueryRowx(`
        INSERT INTO journeys (jou_group_id, jou_passengers, jou_status, jou_car_assigned, jou_creation_date, jou_update_date)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING jou_id
    `, journey.GroupID, journey.Passengers, journey.Status, journey.CarAssignedID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *journeyRepository) GetJourneysByStatus(status enum.JourneyStatus) ([]*entity.Journey, error) {
	var journeys []*model.Journey

	err := r.db.Select(
		&journeys,
		`SELECT *
		FROM journeys j 
		WHERE 
			j.jou_status = $1
			AND j.jou_delete_date IS NULL
		ORDER BY jou_creation_date ASC`,
		status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			var emptyResult []*entity.Journey
			return emptyResult, nil
		}
		return nil, err
	}

	return mapper.ListJourneyModelToEntity(journeys), nil
}

func (r *journeyRepository) CleanJourneys() error {
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM journeys 
		WHERE 1 = 1
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *journeyRepository) UpdateJourneyStatus(journeyID uint, newStatus enum.JourneyStatus) error {
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE journeys 
		SET 
			jou_status = $1,
			jou_update_date = NOW()
		WHERE 
			jou_id = $2
			AND jou_delete_date IS NULL
	`, newStatus, journeyID)

	return tx.Commit()
}

func (r *journeyRepository) GetPendingJourneysWherePassengersLessOrEqualTo(maxPassengers uint) ([]*entity.Journey, error) {
	var journeys []*model.Journey

	err := r.db.Select(
		&journeys,
		`SELECT *
		FROM journeys j 
		WHERE 
			j.jou_status = $1
			AND j.jou_delete_date IS NULL
			AND j.jou_passengers <= $2
		ORDER BY jou_creation_date ASC`,
		enum.JourneyStatusPending, maxPassengers,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			var emptyResult []*entity.Journey
			return emptyResult, nil
		}
		return nil, err
	}

	return mapper.ListJourneyModelToEntity(journeys), nil
}

func (r *journeyRepository) AssignCarToJourney(journeyID uint, carID uint) error {
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE journeys 
		SET 
			jou_status = $1,
			jou_update_date = NOW(),
			jou_car_assigned = $2
		WHERE 
			jou_id = $3
			AND jou_delete_date IS NULL
	`, enum.JourneyStatusAssigned, carID, journeyID)

	return tx.Commit()
}

func (r *journeyRepository) GetJourneysByCarID(carID uint) ([]*entity.Journey, error) {
	var journeys []*model.Journey

	err := r.db.Select(
		&journeys,
		`SELECT *
		FROM journeys j 
		WHERE 
			j.jou_car_assigned = $1
			AND j.jou_delete_date IS NULL
		ORDER BY jou_creation_date ASC`,
		carID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			var emptyResult []*entity.Journey
			return emptyResult, nil
		}
		return nil, err
	}

	return mapper.ListJourneyModelToEntity(journeys), nil
}
