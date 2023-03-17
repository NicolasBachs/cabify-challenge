-- migrate:up
CREATE INDEX car_assigned_idx ON journeys (jou_car_assigned);

-- migrate:down
DROP INDEX car_assigned_idx ON journeys (jou_car_assigned);
