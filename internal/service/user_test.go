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
		input   input
		result  result
		repoDep func(*gomock.Controller) internal.UserRepo
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
		},
		"all is just fine": {
			input: sampleInput,
			repoDep: func(c *gomock.Controller) internal.UserRepo {
				m := mock.NewMockUserRepo(c)
				m.EXPECT().IsEmailInUse(gomock.Eq(sampleInput.email)).Return(false, nil)
				m.EXPECT().Create(gomock.Eq(sampleInput.email), gomock.Any()).Return(nil)
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

			err := user.Create(tt.input.email, tt.input.password)
			assert.Equal(t, tt.result.err, err)
		})
	}
}
