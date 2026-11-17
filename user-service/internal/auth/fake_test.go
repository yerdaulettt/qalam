package auth

import "context"

type testUser struct {
	User
	hash string
}

type fakeRepo struct {
	users map[string]testUser
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: make(map[string]testUser)}
}

func (f *fakeRepo) GetUser(ctx context.Context, username string) (UserVerify, error) {
	u, ok := f.users[username]
	if !ok {
		return UserVerify{}, ErrNotFound
	}

	return UserVerify{Id: u.Id, Role: u.Role, Hash: u.hash}, nil
}

func (f *fakeRepo) GetUsernameStatus(ctx context.Context, username string) (bool, error) {
	u, ok := f.users[username]
	if !ok {
		return true, nil
	}

	if u.Id != 0 {
		return false, nil
	}

	return true, nil
}

func (f *fakeRepo) Register(ctx context.Context, u RegisterReq) (User, error) {
	user := User{Id: len(f.users) + 1, Name: u.Name, Username: u.Username, Role: u.Role}
	tu := testUser{User: user, hash: u.Password}

	f.users[u.Username] = tu

	return user, nil
}

func (f *fakeRepo) GetMyProfile(ctx context.Context, userId int) (User, error) {
	return User{}, nil
}
