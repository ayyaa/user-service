package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"context"
	"errors"
	"database/sql"
	"gopkg.in/guregu/null.v3"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/SawitProRecruitment/UserService/repository"
	gomock "github.com/golang/mock/gomock"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
)

func TestServer_Register(t *testing.T) {
	c := context.Background()
	ctx := middleware.GetEchoContext(c)

	var myRepo repository.RepositoryInterface = &repository.MockRepositoryInterface{}

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		initRouter    func() (echo.Context, *httptest.ResponseRecorder)
		doMock  func(mock *repository.MockRepositoryInterface)
		wantStatus int
	}{
		{
			name: "success",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"fullName": "John Doe", "password": "Tsurayya10*", "phoneNumber": "+62"}`
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().Insert(gomock.Any(), gomock.Any()).
					Return(1, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error insert data",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"fullName": "John Doe", "password": "Tsurayya10*", "phoneNumber": "+62"}`
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().Insert(gomock.Any(), gomock.Any()).
					Return(0, errors.New("errors"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "error unique phone number",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"fullName": "John Doe", "password": "Tsurayya10*", "phoneNumber": "+62"}`
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().Insert(gomock.Any(), gomock.Any()).
					Return(0, errors.New("uq_users_phone"))
			},
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,

			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := repository.NewMockRepositoryInterface(mockCtrl)

			tt.doMock(repo)

			s.Repository = repo

			ctx, rec := tt.initRouter()
			
			s.Register(ctx)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestServer_Login(t *testing.T) {
	c := context.Background()
	ctx := middleware.GetEchoContext(c)

	var myRepo repository.RepositoryInterface = &repository.MockRepositoryInterface{}

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		initRouter    func() (echo.Context, *httptest.ResponseRecorder)
		doMock  func(mock *repository.MockRepositoryInterface)
		wantStatus int
	}{
		{
			name: "error get user by phone",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"password": "Tsurayya10*", "phoneNumber": "+62"}`
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, errors.New("errors"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "error phone number not found",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"password": "Tsurayya10*", "phoneNumber": "+62"}`
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, nil)
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,

			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := repository.NewMockRepositoryInterface(mockCtrl)

			tt.doMock(repo)

			s.Repository = repo

			ctx, rec := tt.initRouter()
			
			s.Login(ctx)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}


func TestServer_GetUser(t *testing.T) {
	c := context.Background()
	ctx := middleware.GetEchoContext(c)

	var myRepo repository.RepositoryInterface = &repository.MockRepositoryInterface{}

	expectedUser := repository.UserRes{
		Id: 57,
		FullName: null.StringFrom("Tsurayya"),
		PhoneNumber: null.StringFrom("+6281219823417"),
	}

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		initRouter    func() (echo.Context, *httptest.ResponseRecorder)
		doMock  func(mock *repository.MockRepositoryInterface)
		wantStatus int
	}{
		{
			name: "error get user by phone",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, errors.New("errors"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "error no rows found",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, sql.ErrNoRows)
			},
			wantStatus: http.StatusNotFound,
		},

		{
			name: "success",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,

			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := repository.NewMockRepositoryInterface(mockCtrl)

			tt.doMock(repo)

			s.Repository = repo

			ctx, rec := tt.initRouter()

			s.GetUser(ctx)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestServer_EditUser(t *testing.T) {
	c := context.Background()
	ctx := middleware.GetEchoContext(c)

	var myRepo repository.RepositoryInterface = &repository.MockRepositoryInterface{}

	expectedUser := repository.UserRes{
		Id: 57,
		FullName: null.StringFrom("Tsurayya"),
		PhoneNumber: null.StringFrom("+6281219823417"),
	}

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		initRouter    func() (echo.Context, *httptest.ResponseRecorder)
		doMock  func(mock *repository.MockRepositoryInterface)
		wantStatus int
	}{
		{
			name: "error get user by phone",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"full_name": "Tsurayya10*", "phone_number": "+62"}`
				req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, errors.New("errors"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "error no rows found",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"full_name": "Tsurayya10*", "phone_number": "+62"}`
				req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, sql.ErrNoRows)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error get phone number",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"full_name": "Tsurayya10*", "phone_number": "+62"}`
				req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, errors.New("errors"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "error phone number",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"full_name": "Tsurayya10*", "phone_number": "+62"}`
				req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
					Return(repository.UserRes{}, errors.New("errors"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "error update phone number already exist",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"full_name": "Tsurayya10*", "phone_number": "+62"}`
				req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("uq_users_phone"))
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "error update data error",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"full_name": "Tsurayya10*", "phone_number": "+62"}`
				req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "success update",
			fields: fields{
				Repository: myRepo,
			},
			args: args{
				ctx: ctx,
			},
			initRouter: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				reqBody := `{"full_name": "Tsurayya10*", "phone_number": "+62"}`
				req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userId", 1)
				return c, rec
			},
			doMock: func(mock *repository.MockRepositoryInterface) {
				mock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).
					Return(expectedUser, nil)
				mock.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,

			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := repository.NewMockRepositoryInterface(mockCtrl)

			tt.doMock(repo)

			s.Repository = repo

			ctx, rec := tt.initRouter()

			s.EditUser(ctx)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
