CREATE TABLE appointment (
    id                  BIGSERIAL PRIMARY KEY,
    patient_id          BIGINT NOT NULL,
    appointment_date    TIMESTAMPTZ NOT NULL
);
