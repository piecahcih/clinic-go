CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(36) NOT NULL,
    last_name  VARCHAR(36) NOT NULL,
    birth_date DATE NOT NULL,
    gender     VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female')),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('patient','doctor')) DEFAULT('patient'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS appointments (
    id           UUID PRIMARY KEY,
    patient_id   UUID NOT NULL REFERENCES users(id),
    doctor_id    UUID NOT NULL REFERENCES users(id),
    start_time   TIMESTAMPTZ NOT NULL,
    description  TEXT,
    status       VARCHAR(20) NOT NULL CHECK (status IN ('booked', 'completed', 'cancelled')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE OR REPLACE FUNCTION appointment_slot(start_time TIMESTAMPTZ)
RETURNS TSTZRANGE AS $$
    SELECT tstzrange(start_time, start_time + interval '15 minutes')
$$ LANGUAGE sql IMMUTABLE;

ALTER TABLE appointments
    ADD CONSTRAINT no_doctor_double_book
    EXCLUDE USING gist (
        doctor_id WITH =,
        appointment_slot(start_time) WITH &&
    );

ALTER TABLE appointments
    ADD CONSTRAINT no_patient_double_book
    EXCLUDE USING gist (
        patient_id WITH =,
        appointment_slot(start_time) WITH &&
    );