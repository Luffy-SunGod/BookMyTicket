package v1

import (
	"fmt"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/common"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/database"
	"github.com/gofiber/fiber/v2"
	"time"
)

func SigninController(ctx *fiber.Ctx) error {

	var User common.UserSignin
	err := ctx.BodyParser(&User)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Unable to parese the request body %v", err))
	}
	fmt.Println("User---->", User)

	db, err := database.ConnecttoDB()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Unable to connect to db , %v", err))
	}

	userId, err := database.UserExist(db, User)
	if err != nil {

		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("User Doesn't exist %v", err))
	}

	correctPasswoord, err := database.CheckPassword(db, User, userId)
	if err != nil || !correctPasswoord {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Err -%v", err))
	}

	if !correctPasswoord {
		return ctx.Status(fiber.StatusBadRequest).SendString("Wrong Password! Try Again!")
	}

	accessToken := common.GenerateAccessToken(userId, User.Email)
	refreshToken := common.GenerateRefreshToken(userId, User.Email)

	saveRefreshTokenInDB, err := database.SaveRefreshTokenInDB(db, userId, refreshToken)
	if err != nil || !saveRefreshTokenInDB {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Unable to dave rtoken to DB!!")
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "accessToken",                  // Cookie name
		Value:    accessToken,                    // Cookie value
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expiration time (24 hours)
		HTTPOnly: true,                           // Prevent client-side JavaScript from accessing the cookie
		Secure:   true,                           // Use cookie only over HTTPS (set to false for HTTP)
	})

	ctx.ClearCookie("user_token")

	ctx.Cookie(&fiber.Cookie{
		Name:     "refreshToken",                 // Cookie name
		Value:    refreshToken,                   // Cookie value
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expiration time (24 hours)
		HTTPOnly: true,                           // Prevent client-side JavaScript from accessing the cookie
		Secure:   true,                           // Use cookie only over HTTPS (set to false for HTTP)
	})

	response := map[string]interface{}{
		"UserId":       userId,
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)

}

func Test(ctx *fiber.Ctx) error {
	return ctx.SendString("json token working properly !!!")
}
