package endpoint

import (
	"database/sql"
	"fmt"
)

func loadDatabaseIntoMemory(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(readStatement)
	if err != nil {
		return nil, fmt.Errorf("Failed to query rows from backup database: %v", err)
	}
	defer rows.Close()

	originSet := make(map[string]bool)
	for rows.Next() {
		var originID string
		err = rows.Scan(&originID)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan rows from backup database: %v", err)
		}
		originSet[originID] = true
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Failed to scan rows from backup database: %v", err)
	}
	return originSet, nil
}
