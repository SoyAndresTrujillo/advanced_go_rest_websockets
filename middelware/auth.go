package middelware

import (
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/helpers"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/server"
	"net/http"
	"strings"
)

// NO_AUTH_NEEDED is a list of routes that don't need authentication
var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

// shouldCheckToken helper's
func shouldCheckToken(route string) bool {
	for _, r := range NO_AUTH_NEEDED {
		if strings.Contains(route, r) {
			return false
		}
	}
	return true
}

// CheckAuthMiddelware is a middelware that checks if user is authenticated
func CheckAuthMiddelware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user is authenticated
			// If not, return 401 status code
			// If yes, call next handler
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			//tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			//_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			//	return []byte(s.Config().JWTSecret), nil
			//})
			//if err != nil {
			//	http.Error(w, err.Error(), http.StatusUnauthorized)
			//	return
			//}

			_, err := helpers.GetJWTAuthorizationInfo(s, w, r)

			if err != nil {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
