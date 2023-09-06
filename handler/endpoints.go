package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"database/sql"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

// (POST /Register)
func (s *Server) Register(ctx echo.Context) error {
	var (
		req generated.RegisterUserRequest
	)

	// Bind data from request
	err := ctx.Bind(&req)
	if err != nil {
		panic(err)
	}

	// Validate struct 
	_, errorMessages := Validation(req)
	if len(errorMessages) > 0 {
		errorString := strings.Join(errorMessages, "\n")
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: errorString,
		})
	}

	// initial mapping add to database
	newUser := repository.UserReq{
		FullName:     req.FullName,
		Password: hashAndSalt([]byte(req.Password)),
		PhoneNumber:    req.PhoneNumber,
	}

	// call method insert to database
	id, err := s.Repository.Insert(ctx.Request().Context(), newUser)
	if err != nil {
		// check if error contains "uq_users_phone", phone number already exist
		if strings.Contains(err.Error(), "uq_users_phone") {
			return ctx.JSON(http.StatusConflict, generated.Message{
				Status:  false,
				Message: ErrorConflictPhoneNumber,
			})
		}

		// return if any error
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// return success
	return ctx.JSON(http.StatusOK, generated.UserId{
		Id: strconv.Itoa(id),
	})
}

// (GET /GetUser)
func (s *Server) GetUser(ctx echo.Context) error {
	id := ctx.Get("userId").(int)

	// Get user by id
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), id)
	if err != nil {
		// check if no rows user
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.Message{
				Status:  false,
				Message: fmt.Sprintf(ErrorUserIdNotFoundf, id),
			})
		}

		// return error
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// build response from database
	resp := generated.UserShort{
		FullName:  user.FullName.ValueOrZero(),
		PhoneNumber: user.PhoneNumber.ValueOrZero(),
	}

	// return success
	return ctx.JSON(http.StatusOK, resp)
}

// (POST /Login)
func (s *Server) Login(ctx echo.Context) error {

	// intial
	var req generated.LoginReq

	// Bind data from request
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// validate if any user by phone
	userData, err := s.Repository.GetUserByPhone(ctx.Request().Context(), req.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// if no sql rows, user not found
	if userData.Id == 0 {
		return ctx.JSON(http.StatusNotFound, generated.Message{
			Status:  false,
			Message: ErrorPhoneNumberNotExist,
		})
	}

	// compare password from request and database
	match := comparePasswords(userData.Password, []byte(req.Password))

	if match {
		// generate token auth
		token, err := generateToken(userData)
		fmt.Println(err)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, ErrorGenerateToken)
		}

		// return succes and user token
		return ctx.JSON(http.StatusOK, generated.SuccessLogin{
			Status:  true,
			Message: SuccesLogin,
			Token:   token,
		})
	}

	// return error
	return ctx.JSON(http.StatusBadRequest, generated.Message{
		Status:  false,
		Message: ErrorUnsuccessfulLogin,
	})

}

// (PATCH /EditUser)
func (s *Server) EditUser(ctx echo.Context) error {
	id := ctx.Get("userId").(int)

	// Bind data from request
	var req generated.UserShort
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// check user by id from db
	userExst, err := s.Repository.GetUserByID(ctx.Request().Context(), id)
	if err != nil && err != sql.ErrNoRows {
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// if no sql rows, user not found
	if userExst.Id == 0 {
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: fmt.Sprintf(ErrorUserIdNotFoundf, id),
		})
	}

	// validate request edit
	_, errorMessages := Validation(req)
	if len(errorMessages) > 0 {
		errorString := strings.Join(errorMessages, "\n")
		return ctx.JSON(http.StatusBadRequest, generated.Message{
			Status:  false,
			Message: errorString,
		})
	}
	// check if phone already used
	userData, err := s.Repository.GetUserByPhone(ctx.Request().Context(), req.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// check if new phone number already user for another user
	if userData.Id != userExst.Id {
		return ctx.JSON(http.StatusConflict, generated.Message{
			Status:  false,
			Message: ErrorConflictPhoneNumber,
		})
	}

	newName := req.FullName
	if req.FullName == "" {
		newName = userExst.FullName.ValueOrZero()
	}

	// mapping data 
	user := repository.UserReq{
		PhoneNumber: req.PhoneNumber,
		FullName:  newName,
	}

	// call method update user
	err = s.Repository.UpdateUser(ctx.Request().Context(), user, id)
	if err != nil {
		// check if error contains "uq_users_phone", phone number already exist
		if strings.Contains(err.Error(), "uq_users_phone") {
			return ctx.JSON(http.StatusConflict, generated.Message{
				Status:  false,
				Message: ErrorConflictPhoneNumber,
			})
		}

		return ctx.JSON(http.StatusForbidden, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
	}

	// return success
	return ctx.JSON(http.StatusOK, generated.Message{
		Status:  true,
		Message: fmt.Sprintf(SuccessUpdateUserById, id),
	})
}
