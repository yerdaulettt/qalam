package auth

import (
	"errors"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	repo := newFakeRepo()
	service := NewAuthService(repo, nil)

	tests := []struct {
		user RegisterReq
		err  error
	}{
		{user: RegisterReq{}, err: ErrEmptyFields},
		{user: RegisterReq{Name: "Test"}, err: ErrEmptyFields},
		{user: RegisterReq{Username: "test"}, err: ErrEmptyFields},
		{user: RegisterReq{Password: "123"}, err: ErrEmptyFields},
		{user: RegisterReq{Name: "Test", Username: "test", Password: "123"}, err: ErrShortPassword},
		{user: RegisterReq{Name: "Test", Username: "test", Password: "test1234", Role: "role"}, err: ErrRole},
		{user: RegisterReq{Name: "Test", Username: "test", Password: "test1234", Role: "user"}, err: nil},
		{user: RegisterReq{Name: "Test", Username: "test", Password: "test1234"}, err: ErrUsername},
		{user: RegisterReq{Name: "Test", Username: "test1", Password: "test1234"}, err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, err := service.Register(t.Context(), tt.user)
			if !errors.Is(err, tt.err) {
				t.Errorf("Error %v, correct %v", err, tt.err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	repo := newFakeRepo()
	jwtAuth := NewJwtAuth([]byte("123"), time.Minute*time.Duration(5), time.Minute*time.Duration(15))
	service := NewAuthService(repo, jwtAuth)

	_, err := service.Register(t.Context(), RegisterReq{Username: "user1", Password: "user1234", Name: "User"})
	if err != nil {
		t.Errorf("Error %v", err)
	}

	tests := []struct {
		user UserLogin
		err  error
	}{
		{user: UserLogin{}, err: ErrEmptyFields},
		{user: UserLogin{Username: "user"}, err: ErrEmptyFields},
		{user: UserLogin{Password: "123"}, err: ErrEmptyFields},
		{user: UserLogin{Username: "user", Password: "123"}, err: ErrNotFound},
		{user: UserLogin{Username: "user1", Password: "123"}, err: ErrIncorrectPassword},
		{user: UserLogin{Username: "user1", Password: "user1234"}, err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, err := service.Login(t.Context(), tt.user)
			if !errors.Is(err, tt.err) {
				t.Errorf("Error %v, correct %v", err, tt.err)
			}
		})
	}
}
