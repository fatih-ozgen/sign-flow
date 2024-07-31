package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() error {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Successfully connected to the database")

	// Create users table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			membership_id CHAR(16) UNIQUE,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Add membership_id column if it doesn't exist
	_, err = db.Exec(`
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='membership_id') THEN
				ALTER TABLE users ADD COLUMN membership_id CHAR(16) UNIQUE;
			END IF;
		END $$;
	`)
	if err != nil {
		return err
	}

	// Update existing records with a membership_id if they don't have one
	rows, err := db.Query("SELECT id FROM users WHERE membership_id IS NULL")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}

		membershipID := generateMembershipID()
		_, err = db.Exec("UPDATE users SET membership_id = $1 WHERE id = $2", membershipID, id)
		if err != nil {
			return err
		}
	}

	// Make membership_id NOT NULL after updating existing records
	_, err = db.Exec(`
		ALTER TABLE users ALTER COLUMN membership_id SET NOT NULL
	`)
	if err != nil {
		return err
	}

	return nil
}

func createUser(membershipID, username, password string) error {
	log.Printf("Attempting to create user: %s with membership ID: %s\n", username, membershipID)
	var result sql.Result
	var err error

	if password == "" {
		log.Println("Creating user without password (OAuth)")
		result, err = db.Exec("INSERT INTO users (membership_id, username, password) VALUES ($1, $2, NULL)", membershipID, username)
	} else {
		log.Println("Creating user with password")
		result, err = db.Exec("INSERT INTO users (membership_id, username, password) VALUES ($1, $2, $3)", membershipID, username, password)
	}

	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return fmt.Errorf("error creating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v\n", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	log.Printf("User created successfully. Rows affected: %d\n", rowsAffected)
	return nil
}

func getUser(usernameOrEmail string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, membership_id, username, password FROM users WHERE username = $1 OR username = $1", usernameOrEmail).Scan(&user.ID, &user.MembershipID, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func getAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT membership_id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.MembershipID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
