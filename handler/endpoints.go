package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

// (GET /Register)
func (s *Server) Register(ctx echo.Context) error {

	var newUser generated.RegisterUserRequest

	err := ctx.Bind(&newUser)

	if err != nil {
		panic(err)
	}

	var resp generated.RegisterUserResponse

	// isValid, errList := UserValidation(newUser)

	// fmt.Println(isValid)
	// fmt.Println(errList)

	// var id = repository.GetTestByIdInput{1}
	id, err := s.Repository.Insert(ctx.Request().Context(), newUser)

	fmt.Println(err)
	fmt.Println(id)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(name)

	return ctx.JSON(http.StatusOK, resp)
}

// func (s *Server) isValidName(name string) (bool, error) {

// 	var err generated.ErrorMessage
// 	if name == "" {
// 		err = "Field Name is Required"
// 		return false, err
// 	}

// 	lenName := len(name)
// 	if lenName <= 3 {
// 		err = "Name must more than 3 character"
// 		return false, err
// 	}

// 	if lenName <= 60 {
// 		err = "Name must less than 60 character"
// 		return false, err
// 	}

// 	return true, nil
// }

// func (s *Server) isValidName(name string) (bool, error) {

// 	var err generated.ErrorMessage
// 	if name == "" {
// 		err = "Field Name is Required"
// 		return false, err
// 	}

// 	lenName := len(name)
// 	if lenName <= 3 {
// 		err = "Name must more than 3 character"
// 		return false, err
// 	}

// 	if lenName <= 60 {
// 		err = "Name must less than 60 character"
// 		return false, err
// 	}

// 	return true, nil
// }

// func (s *Server) Login(ctx echo.Context, params generated.HelloParams) error {

// 	var resp generated.HelloResponse
// 	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
// 	var input = 1
// 	name, err := s.Repository.GetTestById(ctx.Request().Context(), input)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(name)

// 	return ctx.JSON(http.StatusOK, resp)
// }
