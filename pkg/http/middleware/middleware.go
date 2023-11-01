package middleware

import (
	"net/http"
	"strings"

	"github.com/qredo-external/go-rnov/pkg/auth"
	"github.com/qredo-external/go-rnov/pkg/storage"
)

const (
	authHeader = "Authorization"
	basic      = "Basic"
)

// AuthMiddleware - access to auth storage to check token data.
type AuthMiddleware struct {
	auth.Operations
	storage storage.ManageUsers
}

// NewAuthMiddleware - auth middleware constructor.
func NewAuthMiddleware(auth auth.Operations, storage storage.ManageUsers) *AuthMiddleware {
	return &AuthMiddleware{
		Operations: auth,
		storage:    storage,
	}
}

// Authentication - custom HTTP middleware that validates user's basic auth.
func Authentication(auth AuthMiddleware, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jwt string
		ah := r.Header.Get(authHeader)
		jwt, valid := validateAuthStructure(ah)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// note check whether is a valid token - issued by us and still usable -timestamp-
		if err := auth.ValidateJWT(jwt); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// note check whether despite being a valid token it might been invalidated in our system
		if !auth.storage.IsActiveToken(jwt) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// validateAuthStructure - validates that the authorization header value provided by the user has a valid structure.
func validateAuthStructure(ah string) (string, bool) {
	if res := strings.Split(ah, " "); res[0] == basic && len(res) == 2 {
		return res[1], true
	}

	return "", false
}
