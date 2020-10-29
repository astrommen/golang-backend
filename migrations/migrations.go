package migrations

import (
	"golang-backend/helpers"
	"golang-backend/interfaces"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Username string
	Email string
	Password string
}

type Account struct {
	gorm.Model
	Type string
	Name	string
	Balance uint
	UserID uint
}

func createAccounts() {
	db := connectDB()

	users := &[2]interfaces.User {
		{Username: "Anna", Email: "anna@gmail.com"},
		{Username: "Scott", Email: "scott@gmail.com"},
	}

	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := &interfaces.User{
			Username: users[i].Username, 
			Email: users[i].Email, 
			Password: generatedPassword,
		}
		db.Create(&user)

		account := &interfaces.Account{
			Type: "Daily Account", 
			Name: string(users[i].Username + "'s"),
			Balance: uint(10000 * int(i+1)),
			UserID: user.ID,
		}
		db.Create(&account)
	}
	defer db.Close()
}

func Migrate()  {
	User := &interfaces.User{}
	Account := &interfaces.Account{}
	db := connectDB()
	db.AutoMigrate(&User, &Account)
	defer db.Close()

	createAccounts()
}
