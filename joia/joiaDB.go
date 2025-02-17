package joia

import (
	"database/sql"
	"time"
)

var db *sql.DB

type Joia struct {
	ID   int
	PaymentDate time.Time
	EndDate time.Time
}

func CreateJoiaTable(dbV *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS joia (
		id INTEGER PRIMARY KEY,
		payment_date DATETIME,
		end_date DATETIME,
		FOREIGN KEY(id) REFERENCES users(id)
	);
	`
	_, err := dbV.Exec(createTableSQL)
	if err != nil {
		return err
	}

	db = dbV

	return nil
}

func GetLastJoiaById(id int) (Joia, error) {
	now := time.Now()
	query := `
		SELECT id, payment_date, end_date 
		FROM joia 
		WHERE id = ? AND end_date > ? 
		ORDER BY end_date DESC 
		LIMIT 1
	`
	row := db.QueryRow(query, id, now)

	var joia Joia
	err := row.Scan(&joia.ID, &joia.PaymentDate, &joia.EndDate)
	if err != nil {
		return Joia{}, err
	}

	return joia, nil
}

func InsertJoia(joia Joia) error {
	insertJoiaSQL := `
	INSERT INTO joia (id, payment_date, end_date) VALUES (?, ?, ?)
	`
	_, err := db.Exec(insertJoiaSQL, joia.ID, joia.PaymentDate, joia.EndDate)
	if err != nil {
		return err
	}

	return nil
}
