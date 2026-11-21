package fakes

import (
	"context"

	"user-service/internal/user"
)

type testUser struct {
	user.User
	hash string
}

type userRepo struct {
	users map[string]testUser
}

func NewUserRepo() *userRepo {
	return &userRepo{users: make(map[string]testUser)}
}

func (f *userRepo) GetUser(ctx context.Context, username string) (user.UserVerify, error) {
	u, ok := f.users[username]
	if !ok {
		return user.UserVerify{}, user.ErrNotFound
	}

	return user.UserVerify{Id: u.Id, Role: u.Role, Hash: u.hash}, nil
}

func (f *userRepo) GetUsernameStatus(ctx context.Context, username string) (bool, error) {
	u, ok := f.users[username]
	if !ok {
		return true, nil
	}

	if u.Id != 0 {
		return false, nil
	}

	return true, nil
}

func (f *userRepo) Register(ctx context.Context, u user.RegisterReq) (user.User, error) {
	newUser := user.User{Id: len(f.users) + 1, Name: u.Name, Username: u.Username, Role: u.Role}
	tu := testUser{User: newUser, hash: u.Password}

	f.users[u.Username] = tu

	return newUser, nil
}

func (f *userRepo) GetMyProfile(ctx context.Context, userId int) (user.User, error) {
	return user.User{}, nil
}
