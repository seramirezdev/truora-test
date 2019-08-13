package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Obtener la instación de conexión a la base de datos
func GetDB(user, dbname string) (*sql.DB, error) {

	connectionString := fmt.Sprintf("postgresql://%s@localhost:26257/%s?sslmode=disable", user, dbname)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		fmt.Print("Error de conexión a DB", err)
		return nil, err
	}

	createTable := "CREATE TABLE IF NOT EXISTS domains (" +
		"id INT PRIMARY KEY DEFAULT unique_rowid(), " +
		"domain STRING NOT NULL, " +
		"servers JSONB, " +
		"servers_changed BOOL, " +
		"ssl_grade STRING, " +
		"previous_ssl_grade STRING, " +
		"logo STRING, " +
		"title STRING, " +
		"is_down BOOL)"

	if _, err := db.Exec(createTable); err != nil {
		fmt.Println("Error creando la tabla\n", err)
	}

	return db, nil
}
