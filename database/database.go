package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetManagerIDsFromDatabase(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT id FROM managers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managerIDs []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		managerIDs = append(managerIDs, id)
	}

	return managerIDs, nil
}

func GetManagerNameFromDatabase(db *sql.DB, managerID string) (string, error) {
	var name string
	err := db.QueryRow("SELECT name FROM managers WHERE id = ?", managerID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
