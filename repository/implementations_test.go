package repository

import (
	"context"
	"testing"
	"regexp"
	"gopkg.in/guregu/null.v3"
	"database/sql"
	"errors"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	repo := &Repository{DB: mockDB}

	// initial
	user := UserReq{FullName: "Ayya Tsurayya", PhoneNumber: "+6281234590", Password: "pass123"}

	// test case 1 : success
	expectedID := 1
	expectedErr := false

	mock.ExpectQuery(regexp.QuoteMeta(QueryInsertUser)).
		WithArgs(user.FullName, user.PhoneNumber, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	// call method Insert from repo
	id, err := repo.Insert(context.Background(), user)

	assert.Equal(t, id, expectedID)
	assert.Equal(t, expectedErr, err != nil)

	// test case 2 : Error 
	expectedErr = true
	expectedID = 0

	mock.ExpectQuery(regexp.QuoteMeta(QueryInsertUser)).
		WithArgs(user.FullName, user.PhoneNumber, user.Password).
		WillReturnError(errors.New("errors"))

	// call method Insert from repo
	id, err = repo.Insert(context.Background(), user)
	
	assert.Equal(t, id, expectedID)
	assert.Equal(t, expectedErr, err != nil)
}

func TestGetUserByID(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	repo := &Repository{DB: mockDB}

	// intial
	id:= 57

	// test case 1 : success
	expectedUser := UserRes{
		Id: 57,
		FullName: null.StringFrom("Tsurayya"),
		PhoneNumber: null.StringFrom("+6281219823417"),
	}
	expectedErr := false

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"name",
			"phone",
		}).AddRow(
			57,
			"Tsurayya",
			"+6281219823417",
		)

	mock.ExpectQuery(regexp.QuoteMeta(QueryGetUserById)).
		WithArgs(expectedUser.Id).WillReturnRows(rows)

	// call method from repo
	user, err := repo.GetUserByID(context.Background(), id)

	assert.Equal(t, user, expectedUser)
	assert.Equal(t, expectedErr, err != nil)

	// test case 2 : error
	expectedUser = UserRes{}
	expectedErr = true

	mock.ExpectQuery(regexp.QuoteMeta(QueryGetUserById)).
		WithArgs(id).WillReturnError(sql.ErrNoRows)

	// call method from repo
	user, err = repo.GetUserByID(context.Background(), id)

	assert.Equal(t, user, expectedUser)
	assert.Equal(t, expectedErr, err != nil)

}


func TestUpdateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	repo := &Repository{DB: mockDB}

	// initial
	user := UserReq{
		FullName: "Tsurayya",
		PhoneNumber: "+6281219823417",
	}
	id:= 1

	// Test case 1 : success
	expectedErr := false

	mock.ExpectExec(regexp.QuoteMeta(QueryUpdateUser)).
		WithArgs(user.FullName, user.PhoneNumber, id).WillReturnResult(sqlmock.NewResult(1, 1))

	// call method Update from repo
	err := repo.UpdateUser(context.Background(), user, id)
	assert.Equal(t, expectedErr, err != nil)


	// Test case 2 : error
	expectedErr = true

	mock.ExpectExec(regexp.QuoteMeta(QueryUpdateUser)).
		WithArgs(user.FullName, user.PhoneNumber, id).WillReturnError(errors.New("errors"))

	// call method Update from repo
	err = repo.UpdateUser(context.Background(), user, id)
	assert.Equal(t, expectedErr, err != nil)
}

func TestGetUserByPhone(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	repo := &Repository{DB: mockDB}

	// initial
	phone:= "+6281219823417"

	// Test case 1 : success
	expectedUser := UserRes{
		Id: 57,
		FullName: null.StringFrom("Tsurayya"),
		PhoneNumber: null.StringFrom("+6281219823417"),
		Password: "Tsurra12&",
	}
	expectedErr := false
	
	rows := sqlmock.NewRows(
		[]string{
			"id",
			"name",
			"phone",
			"password",
		}).AddRow(
			57,
			"Tsurayya",
			"+6281219823417",
			"Tsurra12&",
		)

	mock.ExpectQuery(regexp.QuoteMeta(QueryGetUserByPhone)).
		WithArgs(phone).WillReturnRows(rows)
	
	// call method GetUserByPhone 
	user, err := repo.GetUserByPhone(context.Background(), phone)

	assert.Equal(t, user, expectedUser)
	assert.Equal(t, expectedErr, err != nil)

	// Test case 2 : Error SQL no rows
	expectedUser = UserRes{}
	expectedErr = false

	mock.ExpectQuery(regexp.QuoteMeta(QueryGetUserByPhone)).
		WithArgs(phone).WillReturnError(sql.ErrNoRows)

	// call method GetUserByPhone 
	user, err = repo.GetUserByPhone(context.Background(), phone)
	assert.Equal(t, user, expectedUser)
	assert.Equal(t, expectedErr, err != nil)

	// Test case 3 : Error
	expectedUser = UserRes{}
	expectedErr = true

	mock.ExpectQuery(regexp.QuoteMeta(QueryGetUserByPhone)).
		WithArgs(phone).WillReturnError(errors.New("errors"))

	// call method GetUserByPhone 
	user, err = repo.GetUserByPhone(context.Background(), phone)
	assert.Equal(t, user, expectedUser)
	assert.Equal(t, expectedErr, err != nil)

}

func TestNewRepository(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	// Create a mock NewRepositoryOptions
	opts := NewRepositoryOptions{
		Dsn: "mock_dsn",
	}

	// Call the NewRepository function with the mock options
	repo := NewRepository(opts)

	// Assert that the Repository's DB field is not nil
	assert.NotNil(t, repo.DB)
}