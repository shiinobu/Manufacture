package helpers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type DB struct {
	DBUSER string
	DBPASS string
	DBNAME string
	DBPORT string
	DBHOST string
}

var RunDB *sql.DB

func InitDB() error {
	if err := godotenv.Load(); err != nil {
		// .env is optional when environment variables are provided by the runtime.
	}

	config := DB{
		DBUSER: os.Getenv("DB_USER"),
		DBPASS: os.Getenv("DB_PASS"),
		DBNAME: os.Getenv("DB_NAME"),
		DBPORT: os.Getenv("DB_PORT"),
		DBHOST: os.Getenv("DB_HOST"),
	}

	if config.DBUSER == "" || config.DBPASS == "" || config.DBNAME == "" || config.DBPORT == "" || config.DBHOST == "" {
		return fmt.Errorf("missing required database environment variables")
	}

	return Connection(config.DBUSER, config.DBPASS, config.DBNAME, config.DBPORT, config.DBHOST)
}

func Connection(dbuser, dbpassword, dbname, dbport, dbhost string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpassword, dbhost, dbport, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	RunDB = db
	return nil
}

func Query(query string) (*sql.Rows, error) {
	return RunDB.Query(query)
}

func QueryParams(query string, params []any) (*sql.Rows, error) {
	return RunDB.Query(query, params...)
}

func QueryExec(query string, params []any) (sql.Result, error) {
	return RunDB.Exec(query, params...)
}
