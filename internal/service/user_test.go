package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/namphung1998/auth-service-go/internal"
	"github.com/namphung1998/auth-service-go/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	type input struct {
		email    string
		password string
	}

	type result struct {
		err error
	}

	sampleInput := input{
		email:    "email",
		password: "password",
	}

	tests := map[string]struct {
		input     input
		result    result
		repoDep   func(*gomock.Controller) internal.UserRepo
		bcryptDep func(*gomock.Controller) internal.BcryptService
	}{
		"email is taken": {
			input: sampleInput,
			result: result{
				err: internal.NewEmailInUseError(sampleInput.email),
			},
			repoDep: func(c *gomock.Controller) internal.UserRepo {
				m := mock.NewMockUserRepo(c)
				m.EXPECT().IsEmailInUse(gomock.Eq(sampleInput.email)).Return(true, nil)
				return m
			},
		},
		"repo.Create returns an error": {
			input: sampleInput,
			result: result{
				err: errors.New("repo.Create"),
			},
			repoDep: func(c *gomock.Controller) internal.UserRepo {
				m := mock.NewMockUserRepo(c)
				m.EXPECT().IsEmailInUse(gomock.Eq(sampleInput.email)).Return(false, nil)
				m.EXPECT().Create(gomock.Eq(sampleInput.email), gomock.Any()).Return(errors.New("repo.Create"))
				return m
			},
			bcryptDep: func(c *gomock.Controller) internal.BcryptService {
				m := mock.NewMockBcryptService(c)
				m.EXPECT().GenerateFromPassword(gomock.Eq([]byte(sampleInput.password))).Return([]byte{}, nil)
				return m
			},
		},
		"all is just fine": {
			input: sampleInput,
			repoDep: func(c *gomock.Controller) internal.UserRepo {
				m := mock.NewMockUserRepo(c)
				m.EXPECT().IsEmailInUse(gomock.Eq(sampleInput.email)).Return(false, nil)
				m.EXPECT().Create(gomock.Eq(sampleInput.email), gomock.Any()).Return(nil)
				return m
			},
			bcryptDep: func(c *gomock.Controller) internal.BcryptService {
				m := mock.NewMockBcryptService(c)
				m.EXPECT().GenerateFromPassword(gomock.Eq([]byte(sampleInput.password))).Return([]byte{}, nil)
				return m
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			user := User{}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			if tt.repoDep != nil {
				user.repo = tt.repoDep(ctrl)
			}

			if tt.bcryptDep != nil {
				user.bcrypt = tt.bcryptDep(ctrl)
			}

			err := user.Create(tt.input.email, tt.input.password)
			assert.Equal(t, tt.result.err, err)
		})
	}
}

func TestLogin(t *testing.T) {
	type input struct {
		email    string
		password string
	}

	type result struct {
		err error
		res internal.LoginResponse
	}

	sampleInput := input{
		email:    "email",
		password: "password",
	}

	sampleRes := internal.LoginResponse{
		Token: "token",
	}

	tests := map[string]struct {
		input     input
		result    result
		repoDep   func(*gomock.Controller) internal.UserRepo
		bcryptDep func(*gomock.Controller) internal.BcryptService
		jwtDep    func(*gomock.Controller) internal.JWTService
	}{
		"all is well": {
			input: sampleInput,
			result: result{
				err: nil,
				res: sampleRes,
			},
			repoDep: func(c *gomock.Controller) internal.UserRepo {
				m := mock.NewMockUserRepo(c)
				m.EXPECT().Get(gomock.Eq(sampleInput.email)).Return(internal.User{}, nil)
				return m
			},
			bcryptDep: func(c *gomock.Controller) internal.BcryptService {
				m := mock.NewMockBcryptService(c)
				m.EXPECT().CompareHashAndPassword(gomock.Eq([]byte{}), gomock.Eq([]byte(sampleInput.password))).Return(nil)
				return m
			},
			jwtDep: func(c *gomock.Controller) internal.JWTService {
				m := mock.NewMockJWTService(c)
				m.EXPECT().GenerateToken(gomock.Eq("")).Return(sampleRes.Token, nil)
				return m
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			user := User{}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			if tt.repoDep != nil {
				user.repo = tt.repoDep(ctrl)
			}

			if tt.bcryptDep != nil {
				user.bcrypt = tt.bcryptDep(ctrl)
			}

			if tt.jwtDep != nil {
				user.jwt = tt.jwtDep(ctrl)
			}

			res, err := user.Login(tt.input.email, tt.input.password)
			assert.Equal(t, tt.result.err, err)
			assert.Equal(t, tt.result.res, res)
		})
	}
}
