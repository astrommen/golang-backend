package users

import (
	"time"

	"golang-backend/helpers"
	"golang-backend/interfaces"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func Login(username string, pass string) map[string]interface{} {
	db := helpers.ConnectDB()
	user := &interfaces.User{}
	passErr := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), 
		[]byte(pass),
	)
	accounts := []interfaces.ResponseAccount{}
	responseUser := &interfaces.ResponseUser{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Accounts: accounts,
	}
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry": time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	var response = map[string]interface{}{
		"message": "all is fine",
	}

	response["jwt"] = token
	response["data"] = responseUser

	if db.Where("username = ?", username).First(&user).RecordNotFound() {
		return map[string]interface{} {
			"message": "User not found",
		}
	}

	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return map[string]interface{} {
			"message": "Wrong password",
		}
	}

	db.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

	defer db.Close()
	
	return response
}