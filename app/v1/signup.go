package v1

import (
	"fmt"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/common"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/database"
	"github.com/gofiber/fiber/v2"
)

func SignupController(ctx *fiber.Ctx) error {
	fmt.Println("signup controller started")
	var User common.UserSignup

	err := ctx.BodyParser(&User)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Unable to parse request body :- %v", err))
	}

	db, err := database.ConnecttoDB()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Unable to connect to db , %v", err))
	}

	userId, err := database.InsertUser(db, User)
	response := map[string]interface{}{
		"userId":  userId,
		"message": "Please redirect to sign in page!!",
	}

	return ctx.Status(fiber.StatusOK).JSON(response)

}
