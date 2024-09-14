package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return fmt.Errorf("DATABASE_URL environment variable is not set")
    }

    var err error
    
    if DB == nil {
        DB, err = sql.Open("postgres", dbURL)
        if err != nil {
            return err
        }
    }

    if err = DB.Ping(); err != nil {
        return err
    }

    return nil
}