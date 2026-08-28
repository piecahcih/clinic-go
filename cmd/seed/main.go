package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

//go run ./cmd/seed

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()

	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer db.Close()

	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("read schema.sql: %v", err)
	}
	for _, stmt := range strings.Split(string(schema), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			log.Fatalf("apply schema: %v", err)
		}
	}

	// TRUNCATE TABLE
	// if _, err := db.ExecContext(ctx, "TRUNCATE TABLE appointments, users"); err != nil {
	// 	log.Fatalf("truncate tables: %v", err)
	// }

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	patientID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	doctorID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	appointmentID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, first_name, last_name, birth_date, gender, email, password_hash, role)
		VALUES
			($1, 'Pichayapa', 'Thaisedhawatkul', '1990-05-14', 'female', 'peach@example.com', $2, 'patient'),
			($3, 'Neerawan', 'Saelee', '1985-11-02', 'male', 'Neer.doctor@example.com', $2, 'doctor')
		ON CONFLICT (id) DO NOTHING`,
		patientID, string(hash), doctorID)
	if err != nil {
		log.Fatalf("seed users: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO appointments (id, patient_id, doctor_id, start_time, description, status)
		VALUES ($1, $2, $3, '2026-09-01T09:00:00Z', 'checkup', 'booked')
		ON CONFLICT (id) DO NOTHING`,
		appointmentID, patientID, doctorID)
	if err != nil {
		log.Fatalf("seed appointments: %v", err)
	}

	fmt.Println("seed complete")
}
