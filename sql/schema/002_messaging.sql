-- +goose Up
CREATE TABLE birthday_notifications (
    id SERIAL PRIMARY KEY,
    person_id varchar(16) NOT NULL,
    scheduled_at varchar(16) NOT NULL,
    state VARCHAR(32) DEFAULT 'scheduled' NOT NULL,
    CONSTRAINT fk_person
        FOREIGN KEY(person_id) 
        REFERENCES persons(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_person_id UNIQUE (person_id)
);

-- +goose Down
DROP TABLE birthday_notifications;