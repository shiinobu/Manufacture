package helpers

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type DB struct {
    DBUSER      string
    DBPASS      string
    DBNAME      string
    DBPORT      string
    DBHOST      string
}

var RunDB *sql.DB

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	DBCONFIG := DB{
		DBUSER: os.Getenv("DB_USER"),
		DBPASS: os.Getenv("DB_PASS"),
		DBNAME: os.Getenv("DB_NAME"),
		DBPORT: os.Getenv("DB_PORT"),
		DBHOST: os.Getenv("DB_HOST"),
	}

	if DBCONFIG.DBUSER == "" || DBCONFIG.DBPASS == "" || DBCONFIG.DBNAME == "" || DBCONFIG.DBPORT == "" || DBCONFIG.DBHOST == "" {
		log.Fatal("Missing required database environment variables")
	}

	return Connection(
		DBCONFIG.DBUSER,
		DBCONFIG.DBPASS,
		DBCONFIG.DBNAME,
		DBCONFIG.DBPORT,
		DBCONFIG.DBHOST,
	)
}

func Connection(dbuser, dbpassword, dbname, dbport, dbhost string) error {
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname
	dbOK, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err = dbOK.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}
	RunDB = dbOK
	return nil
}

func Query(query string) (*sql.Rows, error) {
    rows, err := RunDB.Query(query)
    if err != nil {
        return nil, err
    }
    return rows, nil
}

func QueryParams(query string, params []any) (*sql.Rows, error) {
    rows, err := RunDB.Query(query, params...)
    if err != nil {
        return nil, err
    }
    return rows, nil
}

func QueryExec(query string, params []any) (res sql.Result, err error) {
    res, err = RunDB.Exec(query, params...)
    if err != nil {
        return nil, err
    }
    return res, nil
}