package middleware

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"log"
// 	"net/http"

// 	"github.com/jimmyvallejo/blog-aggregator-go/internal/api/common"
// 	"github.com/jimmyvallejo/blog-aggregator-go/internal/utils"
// )

// func (m *Middleware) IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		extracted, err := utils.ExtractToken(r, "ApiKey ")
// 		if err != nil {
// 			if tokenErr, ok := err.(*utils.TokenError); ok {
// 				http.Error(w, tokenErr.Message, http.StatusBadRequest)
// 			} else {
// 				http.Error(w, "Invalid API key", http.StatusBadRequest)
// 			}
// 			return
// 		}

// 		u, err := m.DB.GetUserByApiKey(r.Context(), extracted)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, sql.ErrNoRows):
// 				http.Error(w, "User not found", http.StatusNotFound)
// 			default:
// 				log.Printf("Error getting user by API key: %v", err)
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), common.UserContextKey, u)

// 		r = r.WithContext(ctx)

// 		next.ServeHTTP(w, r)
// 	}
// }
