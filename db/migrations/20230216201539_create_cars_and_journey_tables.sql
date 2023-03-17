-- migrate:up
CREATE TABLE cars (
  car_id SERIAL PRIMARY KEY,
  car_max_seats INTEGER NOT NULL,
  car_creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  car_update_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  car_delete_date TIMESTAMP
);

CREATE INDEX cars_creation_date_idx ON cars (car_creation_date);
CREATE INDEX cars_delete_date_idx ON cars (car_delete_date);

CREATE TABLE journeys (
  jou_id SERIAL PRIMARY KEY,
  jou_car_assigned INTEGER REFERENCES cars (car_id),
  jou_group_id INTEGER NOT NULL,
  jou_passengers INTEGER NOT NULL,
  jou_status VARCHAR(20) NOT NULL,
  jou_creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  jou_update_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  jou_delete_date TIMESTAMP
);

CREATE INDEX jou_car_assigned_idx ON journeys (jou_car_assigned);
CREATE INDEX jou_group_id_idx ON journeys (jou_group_id);
CREATE INDEX journeys_creation_date_idx ON journeys (jou_creation_date);
CREATE INDEX journeys_delete_date_idx ON journeys (jou_delete_date);

-- migrate:down
DROP TABLE IF EXISTS journeys;
DROP TABLE IF EXISTS cars;