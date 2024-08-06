package middleware

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/common"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/utils"
)

func (m *Middleware) IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Access token not found", http.StatusUnauthorized)
			} else {
				http.Error(w, "Error reading cookie", http.StatusBadRequest)
			}
			return
		}

		validated, err := utils.ValidateToken(cookie.Value, m.JWT)
		if err != nil {
			http.Error(w, "error validating token", http.StatusUnauthorized)
			return
		}

		claims, ok := validated.Claims.(*jwt.RegisteredClaims)
		if !ok {
			http.Error(w, "error getting claims from token", http.StatusInternalServerError)
			return
		}

		idStr, err := claims.GetSubject()
		if err != nil {
			http.Error(w, "error getting subject from claims", http.StatusInternalServerError)
			return
		}

		var user database.User
		err = m.RDB.GetJSON("user"+idStr, &user)
		if err == nil && user.ID != uuid.Nil {
			log.Printf("Cache hit: User retrieved from cache: ID=%s, Email=%s, Handle=%s",
				user.ID, user.Email, user.Handle)

			ctx := context.WithValue(r.Context(), common.UserContextKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			log.Printf("Cache miss: Error retrieving user from cache: %v", err)
		} else {
			log.Printf("Cache miss: User not found in cache or invalid")
		}

		log.Printf("Fetching user from database: ID=%s", idStr)

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid user ID format", http.StatusInternalServerError)
			return
		}

		u, err := m.DB.GetUserByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				http.Error(w, "User not found", http.StatusNotFound)
			default:
				log.Printf("Error getting user by ID: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		log.Printf("User fetched from database: ID=%s, Email=%s, Handle=%s",
			u.ID, u.Email, u.Handle)

		err = m.RDB.SetJson("user"+idStr, u, time.Hour)
		if err != nil {
			log.Printf("Failed to save user to cache: %v", err)
		} else {
			log.Printf("User saved to cache: ID=%s", u.ID)
		}

		ctx := context.WithValue(r.Context(), common.UserContextKey, u)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
