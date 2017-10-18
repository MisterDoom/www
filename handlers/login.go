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

func setToken(w http.ResponseWriter, r *http.Request) {
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

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

	// Redirect the user to his profile
	http.Redirect(w, r, "/", 307)
}

// Middleware to protect private pages
func validate(protectedPage func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// If no Auth cookie is set then return a 404 not found
		cookie, err := req.Cookie("Auth")
		if err != nil {
			http.NotFound(res, req)
			return
		}

		// Return a Token using the cookie
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(res, req)
			return
		}

		// Grab the tokens claims and pass it into the original request
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(req.Context(), "user", *claims)
			protectedPage(res, req.WithContext(ctx))
		} else {
			http.NotFound(res, req)
			return
		}
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates["login"].ExecuteTemplate(w, "layout", reCaptchaSiteKey)

	case "POST":
		setToken(w, r)
	}
}
