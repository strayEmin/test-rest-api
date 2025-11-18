package service

import (
	"context"
	"errors"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/utils/hasher"
	"testing"
)

type authUserRepoMock struct {
	createUserFn func(ctx context.Context, user domain.UserModel) (domain.UserModel, error)
	byEmailFn    func(ctx context.Context, email string) (domain.UserModel, error)
}

func (m *authUserRepoMock) CreateUser(ctx context.Context, user domain.UserModel) (domain.UserModel, error) {
	if m.createUserFn == nil {
		panic("createUserFn is nil")
	}

	return m.createUserFn(ctx, user)
}

func (m *authUserRepoMock) ByEmail(ctx context.Context, email string) (domain.UserModel, error) {
	if m.byEmailFn == nil {
		panic("byEmailFn is nil")
	}

	return m.byEmailFn(ctx, email)
}

type tokenServiceMock struct {
	generateTokenFn func(user domain.User) (string, error)
}

func (m *tokenServiceMock) GenerateToken(user domain.User) (string, error) {
	if m.generateTokenFn == nil {
		panic("generateTokenFn is nil")
	}

	return m.generateTokenFn(user)
}

type loggerAuthMock struct{}

func (m *loggerAuthMock) InitLogger()                                  {}
func (m *loggerAuthMock) Debug(args ...interface{})                    {}
func (m *loggerAuthMock) Debugf(template string, args ...interface{})  {}
func (m *loggerAuthMock) Info(args ...interface{})                     {}
func (m *loggerAuthMock) Infof(template string, args ...interface{})   {}
func (m *loggerAuthMock) Warn(args ...interface{})                     {}
func (m *loggerAuthMock) Warnf(template string, args ...interface{})   {}
func (m *loggerAuthMock) Error(args ...interface{})                    {}
func (m *loggerAuthMock) Errorf(template string, args ...interface{})  {}
func (m *loggerAuthMock) DPanic(args ...interface{})                   {}
func (m *loggerAuthMock) DPanicf(template string, args ...interface{}) {}
func (m *loggerAuthMock) Fatal(args ...interface{})                    {}
func (m *loggerAuthMock) Fatalf(template string, args ...interface{})  {}

func TestAuthService_SignUp_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name := "Emin"
	email := "emin@gmail.com"
	pwd := "12345678"
	generatedToken := "generated-token"
	repoMock := &authUserRepoMock{
		createUserFn: func(ctx context.Context, userModel domain.UserModel) (domain.UserModel, error) {
			userModel.ID = 12
			return userModel, nil
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return generatedToken, nil
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	token, err := svc.SignUp(ctx, name, email, pwd)

	if err != nil {
		t.Fatal(err)
	}
	if token != generatedToken {
		t.Fatal("token should be ", generatedToken)
	}
}

func TestAuthService_SignUp_UserAlreadyExists(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name := "Emin"
	email := "emin@gmail.com"
	pwd := "12345678"
	generatedToken := "generated-token"
	repoMock := &authUserRepoMock{
		createUserFn: func(ctx context.Context, userModel domain.UserModel) (domain.UserModel, error) {
			return domain.UserModel{}, apperrors.ErrUserAlreadyExists
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return generatedToken, nil
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	_, err := svc.SignUp(ctx, name, email, pwd)

	if err == nil {
		t.Fatal("error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrUserAlreadyExists) {
		t.Fatal("expected ErrUSerAlreadyExists, got ", err)
	}
}

func TestAuthService_SignUp_UnexpectedRepoError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name := "Emin"
	email := "emin@gmail.com"
	pwd := "12345678"
	generatedToken := "generated-token"
	repoMock := &authUserRepoMock{
		createUserFn: func(ctx context.Context, userModel domain.UserModel) (domain.UserModel, error) {
			return domain.UserModel{}, errors.New("could not connect to database")
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return generatedToken, nil
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	_, err := svc.SignUp(ctx, name, email, pwd)

	if err == nil {
		t.Fatal("error should be occurred")
	}
}

func TestAuthService_SignUp_TokenGenerationFailed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name := "Emin"
	email := "emin@gmail.com"
	pwd := "12345678"
	repoMock := &authUserRepoMock{
		createUserFn: func(ctx context.Context, userModel domain.UserModel) (domain.UserModel, error) {
			userModel.ID = 12
			return userModel, nil
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return "", errors.New("some error with token generation")
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	_, err := svc.SignUp(ctx, name, email, pwd)

	if err == nil {
		t.Fatal("error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrUnsuccesffulTokenGeneration) {
		t.Fatal("expected ErrUnsuccesffulTokenGeneration, got ", err)
	}
}

func TestAuthService_SignIn_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name := "Emin"
	email := "emin@gmail.com"
	pwd := "12345678"
	hashedPwd, err := hasher.HashPassword(pwd)
	generatedToken := "generated-token"
	repoMock := &authUserRepoMock{
		byEmailFn: func(ctx context.Context, email string) (domain.UserModel, error) {
			return domain.UserModel{
				ID:      123,
				Name:    name,
				Email:   email,
				PwdHash: hashedPwd,
			}, nil
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return generatedToken, nil
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	token, err := svc.SignIn(ctx, email, pwd)

	if err != nil {
		t.Fatal(err)
	}
	if token != generatedToken {
		t.Fatal("token should be ", generatedToken)
	}
}

func TestAuthService_SignIn_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "emin@gmail.com"
	pwd := "12345678"
	generatedToken := "generated-token"
	repoMock := &authUserRepoMock{
		byEmailFn: func(ctx context.Context, email string) (domain.UserModel, error) {
			return domain.UserModel{}, apperrors.ErrUserNotFound
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return generatedToken, nil
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	_, err := svc.SignIn(ctx, email, pwd)

	if err == nil {
		t.Fatal("error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatal("expected ErrUserNotFound, got ", err)
	}
}

func TestAuthService_SignIn_UnexpectedRepoError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "emin@gmail.com"
	pwd := "12345678"
	generatedToken := "generated-token"
	repoMock := &authUserRepoMock{
		byEmailFn: func(ctx context.Context, email string) (domain.UserModel, error) {
			return domain.UserModel{}, errors.New("could not connect to database")
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return generatedToken, nil
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	_, err := svc.SignIn(ctx, email, pwd)

	if err == nil {
		t.Fatal("error should be occurred")
	}
}

func TestAuthService_SignIn_WrongPassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name := "Emin"
	email := "emin@gmail.com"
	pwd := "12345678"
	hashedPwd, err := hasher.HashPassword(pwd)
	wrongPassword := "87654321"
	generatedToken := "generated-token"
	repoMock := &authUserRepoMock{
		byEmailFn: func(ctx context.Context, email string) (domain.UserModel, error) {
			return domain.UserModel{
				ID:      10,
				Name:    name,
				Email:   email,
				PwdHash: hashedPwd,
			}, nil
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return generatedToken, nil
		},
	}

	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	_, err = svc.SignIn(ctx, email, wrongPassword)

	if err == nil {
		t.Fatal("error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrWrongPassword) {
		t.Fatal("expected ErrWrongPassword, got ", err)
	}
}

func TestAuthService_SignIn_TokenGenerationFailed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name := "Emin"
	email := "emin@gmail.com"
	password := "12345678"
	hashedPwd, err := hasher.HashPassword(password)
	repoMock := &authUserRepoMock{
		byEmailFn: func(ctx context.Context, email string) (domain.UserModel, error) {
			return domain.UserModel{
				ID:      123,
				Name:    name,
				Email:   email,
				PwdHash: hashedPwd,
			}, nil
		},
	}
	tokenSvcMock := &tokenServiceMock{
		generateTokenFn: func(u domain.User) (string, error) {
			return "", apperrors.ErrUnsuccesffulTokenGeneration
		},
	}
	svc := NewAuthService(&loggerAuthMock{}, repoMock, tokenSvcMock)

	_, err = svc.SignIn(ctx, email, password)
	if err == nil {
		t.Fatal("error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrUnsuccesffulTokenGeneration) {
		t.Fatal("expected ErrUnsuccesffulTokenGeneration, got ", err)
	}
}
