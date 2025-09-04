-- +goose Up
CREATE TABLE appointments (
    id                 VARCHAR PRIMARY KEY NOT NULL,
    patientId          VARCHAR NOT NULL,
    appointmentDate    TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
);

-- +goose Down
DROP TABLE IF EXISTS appointments;
