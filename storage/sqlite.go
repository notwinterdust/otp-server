package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/notwinterdust/otp-server/models"
)

type DB struct {
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			email         TEXT    NOT NULL UNIQUE,
			password_hash TEXT    NOT NULL
		);
		CREATE TABLE IF NOT EXISTS accounts (
			user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			data    TEXT    NOT NULL
		);
	`)
	return err
}

func (db *DB) CreateUser(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.conn.Exec(
		`INSERT INTO users (email, password_hash) VALUES (?, ?)`,
		email, string(hash),
	)
	return err
}

func (db *DB) AuthenticateUser(email, password string) (*models.User, error) {
	row := db.conn.QueryRow(
		`SELECT id, email, password_hash FROM users WHERE email = ?`, email,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	return &u, nil
}

func (db *DB) GetUserByID(id int64) (*models.User, error) {
	row := db.conn.QueryRow(`SELECT id, email FROM users WHERE id = ?`, id)
	var u models.User
	if err := row.Scan(&u.ID, &u.Email); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &u, nil
}

func (db *DB) GetAccounts(userID int64) ([]models.OTPAccount, error) {
	row := db.conn.QueryRow(`SELECT data FROM accounts WHERE user_id = ?`, userID)
	var raw string
	if err := row.Scan(&raw); err == sql.ErrNoRows {
		return []models.OTPAccount{}, nil
	} else if err != nil {
		return nil, err
	}
	var accounts []models.OTPAccount
	if err := json.Unmarshal([]byte(raw), &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (db *DB) SetAccounts(userID int64, accounts []models.OTPAccount) error {
	raw, err := json.Marshal(accounts)
	if err != nil {
		return err
	}
	_, err = db.conn.Exec(
		`INSERT INTO accounts (user_id, data) VALUES (?, ?)
		 ON CONFLICT(user_id) DO UPDATE SET data = excluded.data`,
		userID, string(raw),
	)
	return err
}

func (db *DB) UserExists() (bool, error) {
	var count int
	err := db.conn.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count > 0, err
}
