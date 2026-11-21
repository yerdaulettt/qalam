package user

import (
	"context"
	"slices"
)

type userRepository interface {
	GetUser(ctx context.Context, username string) (UserVerify, error)
	GetUsernameStatus(ctx context.Context, username string) (bool, error)
	Register(ctx context.Context, u RegisterReq) (User, error)
	GetMyProfile(ctx context.Context, userId int) (User, error)
}

type jwtAuth interface {
	NewTokens(userId int, role string) (JwtTokens, error)
	RefreshAccess(refresh string) (string, error)
}

type UserService struct {
	repo      userRepository
	jwtHelper jwtAuth
}

func NewUserService(r userRepository, j jwtAuth) *UserService {
	return &UserService{repo: r, jwtHelper: j}
}

func (s *UserService) Register(ctx context.Context, u RegisterReq) (User, error) {
	if u.Name == "" || u.Username == "" || u.Password == "" {
		return User{}, ErrEmptyFields
	}

	if len(u.Password) < 8 {
		return User{}, ErrShortPassword
	}

	if u.Role == "" {
		u.Role = "user"
	} else if !slices.Contains([]string{"user", "admin"}, u.Role) {
		return User{}, ErrRole
	}

	ok, err := s.repo.GetUsernameStatus(ctx, u.Username)
	if err != nil {
		return User{}, err
	}

	if !ok {
		return User{}, ErrUsername
	}

	hash, err := hashPassword(u.Password)
	if err != nil {
		return User{}, err
	}

	u.Password = hash
	newUser, err := s.repo.Register(ctx, u)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (s *UserService) Login(ctx context.Context, ul UserLogin) (JwtTokens, error) {
	if ul.Username == "" || ul.Password == "" {
		return JwtTokens{}, ErrEmptyFields
	}

	u, err := s.repo.GetUser(ctx, ul.Username)
	if err != nil {
		return JwtTokens{}, err
	}

	err = checkPassword(u.Hash, ul.Password)
	if err != nil {
		return JwtTokens{}, ErrIncorrectPassword
	}

	tokens, err := s.jwtHelper.NewTokens(u.Id, u.Role)

	return tokens, nil
}

func (s *UserService) TokenRefresh(ctx context.Context, refresh string) (string, error) {
	access, err := s.jwtHelper.RefreshAccess(refresh)
	if err != nil {
		return "", err
	}

	return access, nil
}

func (s *UserService) GetMyProfile(ctx context.Context, userId int) (User, error) {
	u, err := s.repo.GetMyProfile(ctx, userId)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
