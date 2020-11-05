package users

import (
	"time"

	"golang-backend/helpers"
	"golang-backend/interfaces"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	var token = prepareToken(user)
	var response = map[string]interface{}{
		"message": "all is fine",
	}
	response["jwt"] = token
	response["data"] = responseUser

	return response
}

func Login(username string, pass string) map[string]interface{} {
	// Add validation to login
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		},
	)

	if valid {
		// Connect DB
		db := helpers.ConnectDB()
		user := &interfaces.User{}
		if db.Where("username = ? ", username).First(&user).RecordNotFound() {
			return map[string]interface{}{
				"message": "all is fine",
			}
		}

		// Verify Password
		passErr := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(pass),
		)
		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{
				"message": "Wrong password",
			}
		}

		// Find accounts for user
		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

		defer db.Close()

		var response = prepareResponse(user, accounts)

		return response
	} else {
		return map[string]interface{}{
			"message": "not valid values",
		}
	}

}

func Register(username string, email string, pass string) map[string]interface{} {
	// Add validation to registration
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		},
	)

	if valid {

	} else {
		return map[string]interface{}{
			"message": "not valid values",
		}
	}
}
