package database

import (
	"fmt"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/common"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(db *sqlx.DB, user common.UserSignup) (int, error) {
	var userId int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	fmt.Println("hashed password", hashedPassword)
	if err != nil {
		fmt.Println(err)
	}
	err = db.QueryRow(`INSERT INTO Users (username,password,email,gender,number)
							values ($1,$2,$3,$4,$5)
							on conflict (username ,email)
							Do Nothing
							returning userid`, user.UserName, hashedPassword, user.Email, user.Gender, user.Number).Scan(&userId)

	fmt.Println("userId is --->", userId)
	if err != nil {
		return userId, err
	}
	return userId, nil
}
