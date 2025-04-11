package database

import (
	"fmt"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/common"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func UserExist(db *sqlx.DB, user common.UserSignin) (int, error) {
	var userid int
	err := db.QueryRow(`Select userid from users where username=$1 and email=$2`, user.UserName, user.Email).Scan(&userid)
	if err != nil {
		return userid, err
	}
	return userid, nil
}

func CheckPassword(db *sqlx.DB, user common.UserSignin, userId int) (bool, error) {
	var pwd string
	err := db.QueryRow(`select password from users where userid=$1`, userId).Scan(&pwd)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(user.Password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func SaveRefreshTokenInDB(db *sqlx.DB, userid int, refreshToken string) (bool, error) {
	success := false
	err := db.QueryRow(`UPDATE users
SET token = $1
WHERE userid = $2 returning true`, refreshToken, userid).Scan(&success)

	fmt.Println("Suuccess-->", success, refreshToken)
	if err != nil {
		return false, err
	}
	return success, nil

}
