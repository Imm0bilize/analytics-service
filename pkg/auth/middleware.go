package auth

import "net/http"

func (g *GrpcAuth) validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if newAccToken, newRefToken, err := g.RefreshTokens(r.Context(), accessToken.Value, refreshToken.Value); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		} else {
			http.SetCookie(w, &http.Cookie{Name: "access_token", Value: newAccToken, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: newRefToken, Path: "/"})
			next.ServeHTTP(w, r)
		}
	})
}
