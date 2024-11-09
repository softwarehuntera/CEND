package datamanagment



// dropAllTables removes all tables from the database.
func dropAllTables(db *sql.DB) error {
	// Begin a transaction to ensure atomicity.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		// Roll back if there's an error during the transaction.
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Get all table names from the information schema (works for PostgreSQL, MySQL, etc.).
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return fmt.Errorf("failed to query table names: %w", err)
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}

		// Execute a DROP command on each table to delete it.
		if _, err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", tableName, err)
		}
	}

	// Check for errors after iterating through rows.
	if err = rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
	}

	return nil
}