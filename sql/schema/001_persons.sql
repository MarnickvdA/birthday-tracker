-- +goose Up

CREATE TABLE persons (
  id varchar(16) PRIMARY KEY,
  name varchar(128) NOT NULL,
  birth_date varchar(16) NOT NULL
);

-- +goose Down
DROP TABLE persons;