package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/brennosalves/go_restapi/config"
	"github.com/brennosalves/go_restapi/handlers/login"
	"github.com/brennosalves/go_restapi/handlers/util"
)

func SetupRouter(app *fiber.App) {

	// HEARTBEAT
	app.Get("/status", func(c *fiber.Ctx) error {
		currentTime := time.Now()
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Server is running",
			"date":    currentTime.Format("2006-01-02 15:04:05"),
		})
	})

	// LOGIN ENDPOINT
	app.Post("/login", login.PostLogin)

	// TEST AUTHORIZATION
	app.Get("/testauth", authMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "You are authorized",
		})
	})

}

// AUTHORIZATION MIDDLEWARE
func authMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// CHECK IF THE AUTHORIZATION HEADER IS PRESENT
		authString := c.Get("Authorization")
		if authString == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"erros": util.ValidationError{
				Message:  "The token was not informed.",
				Field:    "",
				Location: "header",
				Help:     "Please provide the token: bearer 'token'",
			}})
		}

		// CHECK IF THE AUTHORIZATION HEADER STARTS WITH BEARER
		if !strings.HasPrefix(authString, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"erros": util.ValidationError{
				Message:  "Token format invalid.",
				Field:    "",
				Location: "header",
				Help:     "The authorization header must start with 'Bearer '",
			}})
		}

		// REMOVE BEARER PREFIX FROM THE AUTHORIZATION HEADER
		tokenString := authString[len("Bearer "):]

		// Parse and verify the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// VERIFY THE TOKEN SIGNING METHOD
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// PROVIDE THE SECRET KEY STORED IN THE .ENV FILE USED TO SIGN THE TOKEN
			return []byte(config.ApiJwtSecret), nil
		})
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"erros": util.ValidationError{
				Message:  err.Error(),
				Field:    "",
				Location: "header",
				Help:     "Login to get a new token and try again. If the problem persists, contact the support for assistance.",
			}})
		}

		// SET THE TOKEN CLAIMS AS A VALUE IN THE CONTEXT
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("claims", claims)
		} else {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"erros": util.ValidationError{
				Message:  "Invalid token",
				Field:    "",
				Location: "header",
				Help:     "Login to get a new token and try again.",
			}})
		}

		// CALL THE NEXT MIDDLEWARE OR THE ENDPOINT FUNCTION
		return c.Next()
	}
}
