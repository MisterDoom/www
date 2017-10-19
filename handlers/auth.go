package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Email string `json:"email"`
	// recommended having
	jwt.StandardClaims
}

type User struct {
	Email    string `json:"email"`
	LoggedIn bool   `json:"loggedIn"`
}

func setToken(w http.ResponseWriter, r *http.Request) {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	expireCookie := time.Now().Add(time.Hour * 24)

	email := r.PostFormValue("email")

	secret := []byte("secret")

	claims := Claims{
		email,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:8080",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString(secret)

	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", 307)
}

func authMiddleware(page func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		getUser := func() *User {
			user := &User{
				Email:    "",
				LoggedIn: false,
			}

			cookie, err := r.Cookie("Auth")
			if err != nil {
				return user
			}

			token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected siging method")
				}
				return []byte("secret"), nil
			})
			if err != nil {
				return user
			}

			// Grab the tokens claims and pass it into the original request
			if claims2, ok := token.Claims.(*Claims); ok && token.Valid {
				user = &User{Email: claims2.Email, LoggedIn: true}
			}

			return user
		}

		user := getUser()
		ctx := context.WithValue(r.Context(), "user", user)
		page(w, r.WithContext(ctx))

	})
}

// Middleware to protect private pages
func validate(protectedPage func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*User)

		if user.LoggedIn {
			protectedPage(w, r)
		} else {
			renderTemplate(w, r, "message", messageModel{
				Title:   "Authentication error",
				Message: "You must be authenticated to access this page."})
		}

		return
	})
}
