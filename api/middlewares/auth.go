package middlewares

import (
	"errors"
	"go-fiber-digilib/db/models"
	"go-fiber-digilib/pkg/utils"
	"net/http"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	ID   int         `json:"id"`
	Role models.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey       string
	ExpiresDuration int
}

type contextKey string

const userContextKey = contextKey("user")

// TODO: repair this auth middleware

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(utils.GetConfig("JWT_SECRET_KEY"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

// func (jwtCfg *JWTConfig) GenerateToken(userID int, role models.Role) (string, error) {
// 	expire := jwt.NewNumericDate(time.Now().Local().Add(time.Minute * time.Duration(int64(jwtCfg.ExpiresDuration))))

// 	claims := &JWTCustomClaims{
// 		ID:        userID,
// 		ExpiresAt: expire,
// 		Role:      role,
// 	}

// 	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	token, err := rawToken.SignedString([]byte(jwtCfg.SecretKey))

// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

func GetUser(c fiber.Ctx) (jwt.MapClaims, error) {
	user, ok := c.Locals(userContextKey).(*jwt.Token)
	if !ok || user == nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func GetUserID(c fiber.Ctx) (int, error) {
	claim, err := GetUser(c)

	if err != nil {
		return 0, errors.New("invalid token")
	}

	return claim["id"], nil
}

func VerifyToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)

		if user == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		c.Locals(userContextKey, user)

		userData, err := GetUser(c)
		if userData == nil || err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		c.Locals("userData", userData)

		return c.Next()
	}
}

func VerifyAdmin() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, err := GetUser(c)

		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		if user["role"] != models.Admin {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": "access denied",
			})
		}

		return c.Next()
	}
}

// func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c *echo.Context) error {
// 		user := c.Get("user").(*jwt.Token)

// 		if user == nil {
// 			return c.JSON(http.StatusUnauthorized, map[string]string{
// 				"message": "invalid token",
// 			})
// 		}

// 		ctx := context.WithValue(c.Request().Context(), userContextKey, user)
// 		c.SetRequest(c.Request().WithContext(ctx))

// 		userData, err := GetUser(ctx)
// 		if userData == nil || err != nil {
// 			return c.JSON(http.StatusUnauthorized, map[string]string{
// 				"message": "invalid token",
// 			})
// 		}

// 		c.Set("userData", userData)

// 		return next(c)
// 	}
// }

// func VerifyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c *echo.Context) error {
// 		user, err := GetUser(c.Request().Context())

// 		if err != nil {
// 			return c.JSON(http.StatusUnauthorized, map[string]string{
// 				"message": "invalid token",
// 			})
// 		}

// 		if user.Role != models.Admin {
// 			return c.JSON(http.StatusForbidden, map[string]string{
// 				"message": "access denied",
// 			})
// 		}

// 		return next(c)
// 	}
// }
