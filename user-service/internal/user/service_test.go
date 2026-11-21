package user_test

import (
	"errors"
	"testing"
	"time"

	"user-service/internal/auth"
	"user-service/internal/fakes"
	"user-service/internal/user"
)

func TestRegister(t *testing.T) {
	repo := fakes.NewUserRepo()
	service := user.NewUserService(repo, nil)

	tests := []struct {
		newUser user.RegisterReq
		err     error
	}{
		{newUser: user.RegisterReq{}, err: user.ErrEmptyFields},
		{newUser: user.RegisterReq{Name: "Test"}, err: user.ErrEmptyFields},
		{newUser: user.RegisterReq{Username: "test"}, err: user.ErrEmptyFields},
		{newUser: user.RegisterReq{Password: "123"}, err: user.ErrEmptyFields},
		{newUser: user.RegisterReq{Name: "Test", Username: "test", Password: "123"}, err: user.ErrShortPassword},
		{newUser: user.RegisterReq{Name: "Test", Username: "test", Password: "test1234", Role: "role"}, err: user.ErrRole},
		{newUser: user.RegisterReq{Name: "Test", Username: "test", Password: "test1234", Role: "user"}, err: nil},
		{newUser: user.RegisterReq{Name: "Test", Username: "test", Password: "test1234"}, err: user.ErrUsername},
		{newUser: user.RegisterReq{Name: "Test", Username: "test1", Password: "test1234"}, err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, err := service.Register(t.Context(), tt.newUser)
			if !errors.Is(err, tt.err) {
				t.Errorf("Error %v, correct %v", err, tt.err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	repo := fakes.NewUserRepo()
	jwtAuth := auth.NewJwtAuth([]byte("123"), time.Minute*time.Duration(5), time.Minute*time.Duration(15))
	service := user.NewUserService(repo, jwtAuth)

	_, err := service.Register(t.Context(), user.RegisterReq{Username: "user1", Password: "user1234", Name: "User"})
	if err != nil {
		t.Errorf("Error %v", err)
	}

	tests := []struct {
		login user.UserLogin
		err   error
	}{
		{login: user.UserLogin{}, err: user.ErrEmptyFields},
		{login: user.UserLogin{Username: "user"}, err: user.ErrEmptyFields},
		{login: user.UserLogin{Password: "123"}, err: user.ErrEmptyFields},
		{login: user.UserLogin{Username: "user", Password: "123"}, err: user.ErrNotFound},
		{login: user.UserLogin{Username: "user1", Password: "123"}, err: user.ErrIncorrectPassword},
		{login: user.UserLogin{Username: "user1", Password: "user1234"}, err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, err := service.Login(t.Context(), tt.login)
			if !errors.Is(err, tt.err) {
				t.Errorf("Error %v, correct %v", err, tt.err)
			}
		})
	}
}
