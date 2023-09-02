package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(verifyToken)

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	dbDsn := os.Getenv("DATABASE_URL")

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}

func verifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		publicKey := handler.LoadPublicKey()

		tokenString := ctx.Request().Header.Get("Authorization")
		if tokenString == "" {
			return ctx.JSON(http.StatusUnauthorized, generated.Message{
				Status:  false,
				Message: "Missing token",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &handler.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, generated.Message{
				Status:  false,
				Message: "Invalid token",
			})
		}

		if !token.Valid {
			return ctx.JSON(http.StatusUnauthorized, generated.Message{
				Status:  false,
				Message: "Invalid token",
			})
		}

		claims, ok := token.Claims.(*handler.Claims)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, generated.Message{
				Status:  false,
				Message: "Invalid token claims",
			})
		}

		ctx.Set("phone", claims.Phone)
		return next(ctx)
	}
}
