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
