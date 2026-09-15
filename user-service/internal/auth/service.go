package auth

import "context"

type AuthRepository interface {
	GetUser(ctx context.Context, username string) (UserVerify, error)
	GetUsername(ctx context.Context, username string) (string, error)
	Register(ctx context.Context, u RegisterReq) (User, error)
	GetMyProfile(ctx context.Context, userId int) (User, error)
}

type AuthService struct {
	repo      AuthRepository
	jwtHelper *JwtAuth
}

func NewAuthService(r AuthRepository, j *JwtAuth) *AuthService {
	return &AuthService{repo: r, jwtHelper: j}
}

func (s *AuthService) Register(ctx context.Context, u RegisterReq) (User, error) {
	if u.Name == "" || u.Username == "" || u.Password == "" {
		return User{}, ErrEmptyFields
	}

	if u.Role == "" {
		u.Role = "user"
	}

	uname, err := s.repo.GetUsername(ctx, u.Username)
	if err != nil {
		return User{}, err
	}

	if uname == u.Username {
		return User{}, ErrUsername
	}

	if len(u.Password) < 8 {
		return User{}, ErrShortPassword
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

func (s *AuthService) Login(ctx context.Context, ul UserLogin) (JwtTokens, error) {
	if ul.Username == "" || ul.Password == "" {
		return JwtTokens{}, ErrEmptyFields
	}

	user, err := s.repo.GetUser(ctx, ul.Username)
	if err != nil {
		return JwtTokens{}, err
	}

	err = checkPassword(user.Hash, ul.Password)
	if err != nil {
		return JwtTokens{}, ErrIncorrectPassword
	}

	tokens, err := s.jwtHelper.newTokens(user.Id, user.Role)

	return tokens, nil
}

func (s *AuthService) TokenRefresh(ctx context.Context, refresh string) (string, error) {
	access, err := s.jwtHelper.refreshAccess(refresh)
	if err != nil {
		return "", err
	}

	return access, nil
}

func (s *AuthService) GetMyProfile(ctx context.Context, userId int) (User, error) {
	u, err := s.repo.GetMyProfile(ctx, userId)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
