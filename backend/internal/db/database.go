package db

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) (*gorm.DB, error) {
    if dsn == "" {
        dsn = "root:@tcp(localhost:3306)/reward_system?charset=utf8mb4&parseTime=True&loc=Local"
    }

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

	DB = db
	log.Println("Database connected successfully")
	return db, nil
}

func Migrate(db *gorm.DB) error { return nil }