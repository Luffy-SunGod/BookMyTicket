package common

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strings"
	"time"
)

func JWTMiddleware(ctx *fiber.Ctx) error {

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).SendString("No authorization header in the request")
	}
	tokenString := strings.Split(authHeader, "Bearer ")[1]
	log.Println("Token: ", tokenString)

	token, err := jwt.NewParser().Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ctx.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Unexpected signing method: %s", t.Header["alg"]))
		}
		return []byte("access-token-secret-string"), nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return ctx.Status(fiber.StatusUnauthorized).SendString(fmt.Sprint("Error: Claim time isnt correct"))
		}
		body := make(map[string]interface{})
		if err := ctx.BodyParser(&body); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Error: Failed to decode request body")
		}
		fmt.Println("claims----->", ok)
		//sub, _ := claims["user_id"]
		//body["user_id"] = sub

		newBody, err := json.Marshal(body)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to re-encode request body"})
		}
		fmt.Println("NewBody!!", newBody)
		ctx.Request().SetBody(newBody)
		return ctx.Next()
	} else {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Error: Claim isnt correct ")
	}

}

type MyClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID int, email string) string {
	claims := MyClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "myApp",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte("access-token-secret-string"))

	return tokenString

}

func GenerateRefreshToken(userID int, email string) string {
	claims := MyClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "myApp",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // Refresh token expires in 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte("access-token-secret-string"))

	return tokenString

}

//
//func RefreshAccessToken(ctx *fiber.Ctx) error {
//
//}

//Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJlbWFpbCI6ImpvaG5kb2VAZXhhbXBsZS5jb20iLCJpc3MiOiJteUFwcCIsImV4cCI6MTc0MzA3NDc0NywiaWF0IjoxNzQzMDczODQ3fQ.V05QwWCCpNJQWSc9TctfQq2ATeBNOwpfFB9Ti_c4IYk
