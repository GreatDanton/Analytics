package sessions

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/greatdanton/analytics/src/global"
)

// User struct that holds user data -> creating token with this data
type User struct {
	ID       string
	Username string
	LoggedIn bool
}

// CreateUserSession creates session for current user, used for loggin user in.
// Currently Logged In user is defined with Cookie that contains jwt string
// with the relevant user data (id: user.id, username: user.username, loggedIn: bool)
func CreateUserSession(w http.ResponseWriter, userID, username string) error {
	cookie, err := createCookie(userID, username)
	if err != nil {
		return err
	}
	// set cookie that that defines logged in user in users browser
	http.SetCookie(w, &cookie)
	return nil
}

// DestroyUserSession destroys session by changing/emptying fields in currently active
// user cookie
func DestroyUserSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := destroyCookie(r)
	if err != nil {
		// no cookie is present but the user press logout
		// TODO: how do we deal with this
		if err == http.ErrNoCookie {
			log.Println("destroySession, no cookie is present but destroy is called:", err)
			return err
		}
		log.Println(err)
		return err
	}
	http.SetCookie(w, &cookie)
	return nil
}

// CreateToken creates signed token string with user id and username as payload
// signed string password is parsed from config.json file
func createToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = user.ID //User{ID: user.ID, Username: user.Username}
	claims["username"] = user.Username
	claims["loggedIn"] = user.LoggedIn
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims
	tokenEncode := []byte(global.JwtTokenPassword)

	signedString, err := token.SignedString(tokenEncode)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

// CreateCookie creates cookie out of inserted arguments (user.id, user.username)
// and returns cookie or error when that is not possible
func createCookie(id, username string) (http.Cookie, error) {
	// set cookie and redirect
	expiration := time.Now().Add(7 * 24 * time.Hour) // cookie expires in 1 week
	u := User{ID: id, Username: username, LoggedIn: true}
	tokenString, err := createToken(u)
	if err != nil {
		e := fmt.Errorf("CreateCookie: CreateToken:%v", err)
		return http.Cookie{}, e
	}
	cookie := http.Cookie{Name: "Analytics", Value: tokenString,
		Expires: expiration, Path: "/", HttpOnly: true}

	return cookie, err
}

// DestroyCookie since we can't delete cookie on all browsers,
// it sets value of authentication cookie to blank and add expiration date = now
func destroyCookie(r *http.Request) (http.Cookie, error) {
	_, err := r.Cookie("Analytics")
	// cookie does not exist
	if err != nil {
		return http.Cookie{}, err
	}
	c := http.Cookie{Name: "Analytics", Value: "", Expires: time.Now(), Path: "/", HttpOnly: true}
	return c, nil
}
