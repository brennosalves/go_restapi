package login

import (
	"net/http"
	"time"

	"github.com/brennosalves/go_restapi/config"
	"github.com/brennosalves/go_restapi/handlers/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// STRUCT LOGIN
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TOKEN STRUCTS
type TokenPayload struct {
	Username string `json:"username"`
}
type Token struct {
	Token     string `json:"accessToken"`
	Username  string `json:"username"`
	ExpiresIn string `json:"expiresIn"`
}

// FUNCTION FOR POST LOGIN
func PostLogin(c *fiber.Ctx) error {

	var erros []util.ValidationError
	var loginReq Login

	// PARSE THE BODY REQUEST AND CHECK IF IT IS VALID
	if err := c.BodyParser(&loginReq); err != nil {
		erros = append(erros, util.ValidationError{
			Message:  "Invalid request",
			Field:    "",
			Location: "body",
			Help:     "Please check if the body structure is correct.",
		})
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erros": erros})
	}

	// MANUAL CHECK TO GRANT ACCESS
	if loginReq.Username == "username" && loginReq.Password == "password" {

		// GENERATE TOKEN
		token, err := generateToken(TokenPayload{Username: loginReq.Username})
		if err != nil {
			erros = append(erros, util.ValidationError{
				Message:  err.Error(),
				Field:    "",
				Location: "handler",
				Help:     "Fail to generate token. Contact the support.",
			})
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"erros": erros})
		}

		// RETURNS THE TOKEN
		return c.Status(http.StatusOK).JSON(token)
	} else {
		// BAD CREDENTIALS
		erros = append(erros, util.ValidationError{
			Message:  "Invalid Username and/or password.",
			Field:    "Username/password",
			Location: "body",
			Help:     "Check your Username and password and try again.",
		})
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"erros": erros})
	}
}

// FUNCTION TO GENERATE JWT TOKEN
func generateToken(payload TokenPayload) (Token, error) {

	// CREATES A NEW JWT TOKEN
	token := jwt.New(jwt.SigningMethodHS256)

	// SET CLAIMS FOR THE TOKEN
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = payload.Username
	claims["expiresIn"] = time.Now().Add(time.Hour * 24).Unix()

	// SIGN THE TOKEN WITH A SECRET KEY STORE IN THE .ENV FILE
	tokenString, err := token.SignedString([]byte(config.ApiJwtSecret))

	// BUILDS THE TOKEN RESPONSE AND RETURNS
	var t Token
	t.Token = tokenString
	t.Username = payload.Username
	t.ExpiresIn = time.Unix(time.Now().Add(time.Hour*24).Unix(), 0).Format("2006-01-02 15:04:05")

	if err != nil {
		return t, err
	}

	return t, nil
}
