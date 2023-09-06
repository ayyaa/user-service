package main

import (
	"log"
	"os"
	"context"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func main() {

	e := echo.New()

	swagger, err := generated.GetSwagger()

	if err != nil {
		os.Exit(1)
	}

	validatorOptions := &middleware.Options{}

	validatorOptions.Options.AuthenticationFunc = func(c context.Context, input *openapi3filter.AuthenticationInput) error {
		return handler.MiddlewareAuth(c, input)
	}

	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))
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