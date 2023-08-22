package biz

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"kratos-realworld/internal/conf"
	"kratos-realworld/internal/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID       uint
	Username string
	Token    string
	Email    string

	PasswordHash string
}
type UserLogin struct {
	Email    string
	Username string

	Token string
}
type UserUpdate struct {
	Email    string
	Username string
	Password string
}
type ProfileRepo interface {
	GetProfile(ctx context.Context, username string) (*Profile, error)
	FollowUser(ctx context.Context, currentUserID uint, followingID uint) error
	UnfollowUser(ctx context.Context, currentUserID uint, followingID uint) error
	GetUserFollowingStatus(ctx context.Context, currentUserID uint, userIDs []uint) (following []bool, err error)
}
type Profile struct {
	ID        uint
	Username  string
	Following bool
}

func hashPassword(pwd string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func verifyPassword(hashed, input string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input)); err != nil {
		return false
	}
	return true
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
}

type UserUsecase struct {
	ur   UserRepo
	pr   ProfileRepo
	jwtc *conf.JWT
	tm   Transaction
	log  *log.Helper
}

func NewUserUsecase(ur UserRepo, logger log.Logger, jwtc *conf.JWT, tm Transaction) *UserUsecase {
	return &UserUsecase{ur: ur, jwtc: jwtc, tm: tm, log: log.NewHelper(logger)}
}
func (uc *UserUsecase) generateToken(userID uint) string {
	return auth.GenerateToken(uc.jwtc.Secret, userID)
}
func (uc *UserUsecase) Register(ctx context.Context, username, email, password string) (*UserLogin, error) {

	u := &User{
		Email:        email,
		Username:     username,
		PasswordHash: hashPassword(password),
	}
	if err := uc.ur.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return &UserLogin{
		Email:    email,
		Username: username,
		Token:    uc.generateToken(u.ID),
	}, nil

}

func (uc *UserUsecase) Login(ctx context.Context, email, password string) (*UserLogin, error) {
	u, err := uc.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !verifyPassword(u.PasswordHash, password) {
		return nil, errors.Unauthorized("user", "login failed")
	}
	return &UserLogin{
		Email:    u.Email,
		Username: u.Username,

		Token: uc.generateToken(u.ID),
	}, nil
}
func (uc *UserUsecase) GetCurrentUser(ctx context.Context) (*User, error) {
	cu := auth.FromContext(ctx)
	if cu == nil {
		panic("cu is nil")
	}
	u, err := uc.ur.GetUserByID(ctx, cu.UserID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, uu *UserUpdate) (*UserLogin, error) {
	cu := auth.FromContext(ctx)
	u, err := uc.ur.GetUserByID(ctx, cu.UserID)
	if err != nil {
		return nil, err
	}
	u.Email = uu.Email
	u.PasswordHash = hashPassword(uu.Password)
	u, err = uc.ur.UpdateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return &UserLogin{
		Email:    u.Email,
		Username: u.Username,
		Token:    uc.generateToken(u.ID),
	}, nil
}
func (uc *ProfileUsecase) GetProfile(ctx context.Context, username string) (rv *Profile, err error) {
	return uc.pr.GetProfile(ctx, username)
}
