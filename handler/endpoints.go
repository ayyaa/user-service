package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

// (GET /Register)
func (s *Server) Register(ctx echo.Context) error {
	var (
		req generated.RegisterUserRequest
	)

	err := ctx.Bind(&req)
	if err != nil {
		panic(err)
	}

	_, errorMessages := UserValidation(req)
	if len(errorMessages) > 0 {
		errorString := strings.Join(errorMessages, "\n")
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: errorString,
		})
	}

	newUser := repository.UserReq{
		Name:     req.Name,
		Password: hashAndSalt([]byte(req.Password)),
		Phone:    req.Phone,
	}
	id, err := s.Repository.Insert(ctx.Request().Context(), newUser)
	if err != nil {
		if strings.Contains(err.Error(), "uq_users_phone") {
			return ctx.JSON(http.StatusConflict, generated.Message{
				Status:  false,
				Message: "Conflict phone number already exist",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.UserId{
		Id: strconv.Itoa(id),
	})
}

func (s *Server) GetUser(ctx echo.Context) error {

	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: "Error Converting string to int",
		})
	}
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), userId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ctx.JSON(http.StatusNotFound, generated.Message{
				Status:  false,
				Message: "User id not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	resp := generated.UserShort{
		Name:  user.Name.ValueOrZero(),
		Phone: user.Phone.ValueOrZero(),
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {

	var req generated.LoginReq

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	userData, err := s.Repository.GetUserByPhone(ctx.Request().Context(), req.Phone)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	if userData.Id == 0 {
		return ctx.JSON(http.StatusNotFound, generated.Message{
			Status:  false,
			Message: "Phone number not found, please register",
		})
	}

	match := comparePasswords(userData.Password, []byte(req.Password))

	if match {
		token, err := generateToken(userData.Phone.ValueOrZero())
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "Error generating token")
		}

		return ctx.JSON(http.StatusOK, generated.SuccessLogin{
			Status:  true,
			Message: fmt.Sprintf("Success Login user"),
			Token:   token,
		})
	}

	return ctx.JSON(http.StatusUnauthorized, generated.Message{
		Status:  false,
		Message: "Unauthorized",
	})

}

func (s *Server) EditUser(ctx echo.Context) error {

	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	userExst, err := s.Repository.GetUserByID(ctx.Request().Context(), userId)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	if userExst.Id == 0 {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: fmt.Sprintf("User id %d not found", userId),
		})
	}

	var req generated.UserShort

	err = ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	_, errorMessages := UserValidation(req)
	if len(errorMessages) > 0 {
		errorString := strings.Join(errorMessages, "\n")
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: errorString,
		})
	}

	userData, err := s.Repository.GetUserByPhone(ctx.Request().Context(), req.Phone)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	if userData.Id != userExst.Id {
		return ctx.JSON(http.StatusConflict, generated.Message{
			Status:  false,
			Message: "Conflict phone number already exist",
		})
	}

	newName := req.Name
	if req.Name == "" {
		newName = userExst.Name.ValueOrZero()
	}

	user := repository.UserReq{
		Phone: req.Phone,
		Name:  newName,
	}

	err = s.Repository.EditUser(ctx.Request().Context(), user, userId)
	if err != nil {
		if strings.Contains(err.Error(), "uq_users_phone") {
			return ctx.JSON(http.StatusConflict, generated.Message{
				Status:  false,
				Message: "Conflict phone number already exist",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.Message{
		Status:  true,
		Message: fmt.Sprintf("Success update user id %d", userId),
	})
}

func (s *Server) Protected(ctx echo.Context) error {

	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: "Error Converting string to int",
		})
	}
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), userId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ctx.JSON(http.StatusNotFound, generated.Message{
				Status:  false,
				Message: "User id not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	resp := generated.UserShort{
		Name:  user.Name.ValueOrZero(),
		Phone: user.Phone.ValueOrZero(),
	}

	return ctx.JSON(http.StatusOK, resp)
}
