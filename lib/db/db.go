package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func tableExists(db *sql.DB, tableName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createTableIfNotExists(db *sql.DB, tableName string, schema string) error {
	exists, err := tableExists(db, tableName)
	if err != nil {
		return err
	}
	if !exists {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
		log.Printf("Created Table %s\n", tableName)
	}
	return nil
}

func Connect(timeout time.Duration, dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./"+dbname+".db")
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB\n", err)
		return nil, err
	}

	log.Printf("Connected to DB %s successfully\n", dbname)

	tables := []struct {
		name   string
		schema string
	}{
		{usersTable, usersTableSchema},
		{userAssetsTable, userAssetsTableSchema},
		{userTokensTable, userTokensTableSchema},
	}

	for _, table := range tables {
		err = createTableIfNotExists(db, table.name, table.schema)
		if err != nil {
			log.Printf("Error creating table %s: %s\n", table.name, err)
			return nil, err
		}
	}

	return db, nil
}
