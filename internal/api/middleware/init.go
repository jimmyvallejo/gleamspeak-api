package middleware

import "github.com/jimmyvallejo/gleamspeak-api/internal/database"

type Middleware struct {
	DB *database.Queries
}

func NewMiddleware(db *database.Queries) *Middleware {
	return &Middleware{
		DB: db,
	}
}
