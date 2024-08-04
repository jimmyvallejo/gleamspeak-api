package main

import "github.com/jimmyvallejo/gleamspeak-api/internal/database"

type APIConfig struct {
	Port string
	DB   *database.Queries
}
