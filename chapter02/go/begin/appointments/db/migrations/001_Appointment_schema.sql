-- +goose Up
CREATE TABLE appointments (
    id                 VARCHAR PRIMARY KEY NOT NULL,
    patientId          VARCHAR NOT NULL,
    appointmentDate    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS appointments;
