package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/greatdanton/analytics/src/global"
)

// LoginUser is used to gather data from database query
type LoginUser struct {
	ID           string
	Username     string
	PasswordHash []byte
}

// ErrUserDoesNotExist error message if user does not exist in database
var ErrUserDoesNotExist = errors.New("This user does not exist")

// GetUserLoginData gets user data from the database and returns
// User and
func GetUserLoginData(username string, password string) (LoginUser, error) {
	var (
		id       string
		user     string
		passHash []byte
	)

	tryingUser := LoginUser{}
	err := global.DB.QueryRow(`Select id, username, password from users
								where username=$1`, username).Scan(&id, &user, &passHash)
	if err != nil {
		// if no rows are returned, user does not exist
		if err == sql.ErrNoRows {
			return tryingUser, ErrUserDoesNotExist
		}
		return tryingUser, fmt.Errorf("GetUserLoginData: A database error occured: %v", err)
	}
	tryingUser.ID = id
	tryingUser.Username = username
	tryingUser.PasswordHash = passHash

	return tryingUser, nil
}
