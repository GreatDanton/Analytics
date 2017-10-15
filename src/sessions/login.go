package sessions

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/greatdanton/analytics/src/global"
)

// LoggedInUser checks if Analytics (authentication) cookie is present
// in client request and returns user struct with
// {ID: user.id, Username: user.username, LoggedIn: true/false}
func LoggedInUser(r *http.Request) User {
	cookie, err := r.Cookie("Analytics")
	if err != nil {
		return User{}
	}
	tokenString := cookie.Value
	u, loggedIn := GetUserData(tokenString)
	if !loggedIn {
		return User{}
	}
	return u
}

// GetUserData gets userData from JWT token string and returns
// UserData (users.id, users.username)
// loggedIn (bool): true if token is formed correctly
//				  false if token is forged or an error occured
func GetUserData(tokenString string) (User, bool) {
	tokenEncode := []byte(global.JwtTokenPassword)

	claims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenEncode, nil
	})
	u := User{}

	// if error occured
	if err != nil {
		return u, false
	}

	// check if token is valid
	if !token.Valid {
		return u, false
	}

	id := token.Claims.(jwt.MapClaims)["id"]
	username := token.Claims.(jwt.MapClaims)["username"]
	loggedIn := token.Claims.(jwt.MapClaims)["loggedIn"]

	// type assertion, checking if returned values are actually strings
	// if they are not return empty user struct
	if idStr, ok := id.(string); ok {
		u.ID = idStr
	} else {
		return u, false
	}
	if usernameStr, ok := username.(string); ok {
		u.Username = usernameStr
	} else {
		return u, false
	}
	if loggedInStr, ok := loggedIn.(bool); ok {
		u.LoggedIn = loggedInStr
	} else {
		return u, false
	}

	return u, true
}
