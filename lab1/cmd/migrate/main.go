package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Reaction{},
		&ds.Synthesis{},
		&ds.SynthesisReaction{},
		&ds.Users{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
