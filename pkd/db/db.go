package db

import (
	"database/sql"
	"fmt"
	"os"

	//"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema string = `
        CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date CHAR(8) NOT NULL DEFAULT "",
        title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
    );`

func OpenDB(dbFile string) error {
	if dbFile == "" {
		dbFile = "./scheduler.db"
	}
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("database opening error %w", err)
	}
	return nil
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func Init(dbFile string) error {
	if dbFile == "" {
		dbFile = "./scheduler.db"
	}
	//проверяем существование файла
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	if err != nil && !install {
		return fmt.Errorf("cannot check DB file: %w", err)
	}

	//открываем соединение с бд и автоматичеси сздаем файл
	err = OpenDB(dbFile)
	if err != nil {
		return fmt.Errorf("database opening error %w", err)
	}
	if err = db.Ping(); err != nil {
		return fmt.Errorf("database connection error: %w", err)

	}
	// Если БД новая, создаем таблицы
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("error when sending the creation request %w", err)
		}
	}
	return nil
}
