package app

import (
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/common"
	v1 "github.com/Luffy-SunGod/MyShowSeat/cmd/app/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB *sqlx.DB
}

func (c *Service) Start() {
	mux := fiber.New()
	mux.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Accept,Authorization,Content-Type,X-CSRF-TOKEN",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))
	//checking for ping whether our server is up or not
	mux.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK,server is working my friend !!!")
	})

	mux.Post("/signup", v1.SignupController)
	mux.Post("/signin", v1.SigninController)
	mux.Use(common.JWTMiddleware).Post("/createShow", v1.CreateShowController)

	mux.Get("/isSeatFull", v1.IsSeatFullController).Use(common.JWTMiddleware)

	mux.Listen(":8000")

}
