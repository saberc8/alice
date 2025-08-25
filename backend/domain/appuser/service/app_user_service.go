package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	appentity "alice/domain/appuser/entity"
	apprepo "alice/domain/appuser/repository"
	"alice/infra/config"
)

var (
	ErrAppUserNotFound       = errors.New("app user not found")
	ErrAppUserExists         = errors.New("app user already exists")
	ErrAppInvalidCredentials = errors.New("invalid credentials")
	ErrAppUserInactive       = errors.New("user is inactive")
)

type AppUserService interface {
	Register(email, password, nickname string) (*appentity.AppUser, error)
	Login(email, password string) (string, error)
	GetByID(id uint) (*appentity.AppUser, error)
	UpdateProfile(id uint, nickname, avatar, gender, bio string) (*appentity.AppUser, error)
	GetByIDs(ids []uint) ([]*appentity.AppUser, error)
}

type appUserServiceImpl struct {
	repo apprepo.AppUserRepository
}

func NewAppUserService(repo apprepo.AppUserRepository) AppUserService {
	return &appUserServiceImpl{repo: repo}
}

func (s *appUserServiceImpl) Register(email, password, nickname string) (*appentity.AppUser, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if u, _ := s.repo.GetByEmail(email); u != nil {
		return nil, ErrAppUserExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &appentity.AppUser{Email: email, PasswordHash: string(hash), Nickname: nickname, Status: appentity.AppUserStatusActive}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *appUserServiceImpl) Login(email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	u, err := s.repo.GetByEmail(email)
	if err != nil || u == nil {
		return "", ErrAppInvalidCredentials
	}
	if !u.IsActive() {
		return "", ErrAppUserInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrAppInvalidCredentials
	}
	return s.generateToken(u.ID)
}

func (s *appUserServiceImpl) GetByID(id uint) (*appentity.AppUser, error) {
	u, err := s.repo.GetByID(id)
	if err != nil || u == nil {
		return nil, ErrAppUserNotFound
	}
	return u, nil
}

func (s *appUserServiceImpl) UpdateProfile(id uint, nickname, avatar, gender, bio string) (*appentity.AppUser, error) {
	u, err := s.repo.GetByID(id)
	if err != nil || u == nil {
		return nil, ErrAppUserNotFound
	}
	if nickname != "" {
		u.Nickname = nickname
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	if gender != "" { // 简单校验：限制枚举
		g := strings.ToLower(gender)
		if g == "male" || g == "female" || g == "other" { // 允许值
			u.Gender = g
		}
	}
	if bio != "" {
		u.Bio = bio
	}
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *appUserServiceImpl) GetByIDs(ids []uint) ([]*appentity.AppUser, error) {
	return s.repo.GetByIDs(ids)
}

func (s *appUserServiceImpl) generateToken(userID uint) (string, error) {
	cfg := config.Load()
	claims := jwt.MapClaims{
		"app_user_id": userID,
		"exp":         time.Now().Add(time.Duration(cfg.JWT.ExpiresIn) * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.SecretKey))
}
