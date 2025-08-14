package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"alice/domain/user/entity"
	"alice/domain/user/repository"
	"alice/infra/config"
	"alice/pkg/logger"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
)

// UserService 用户领域服务接口
type UserService interface {
	// Register 用户注册
	Register(username, password, email string) (*entity.User, error)

	// Login 用户登录
	Login(username, password string) (string, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(userID uint) (*entity.User, error)

	// UpdateProfile 更新用户资料
	UpdateProfile(userID uint, email string) (*entity.User, error)

	// CreateUserByAdmin 管理员创建用户
	CreateUserByAdmin(username, password, email string, status entity.UserStatus) (*entity.User, error)

	// GetUser 管理员根据ID获取用户
	GetUser(userID uint) (*entity.User, error)

	// ListUsers 分页列出用户
	ListUsers(page, pageSize int) ([]*entity.User, int64, error)

	// UpdateUserByAdmin 管理员更新用户基本信息
	UpdateUserByAdmin(userID uint, email string, status entity.UserStatus, password *string) (*entity.User, error)

	// DeleteUser 管理员删除用户
	DeleteUser(userID uint) error
}

// userServiceImpl 用户服务实现
type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *userServiceImpl) Register(username, password, email string) (*entity.User, error) {
	// 检查用户是否已存在
	existingUser, _ := s.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在
	existingUser, _ = s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		return nil, err
	}

	// 创建用户
	user := &entity.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		Status:       entity.UserStatusActive,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		logger.Errorf("Failed to create user: %v", err)
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *userServiceImpl) Login(username, password string) (string, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != entity.UserStatusActive {
		return "", ErrUserInactive
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// 生成JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		logger.Errorf("Failed to generate token: %v", err)
		return "", err
	}

	return token, nil
}

// GetUserByID 根据ID获取用户
func (s *userServiceImpl) GetUserByID(userID uint) (*entity.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// UpdateProfile 更新用户资料
func (s *userServiceImpl) UpdateProfile(userID uint, email string) (*entity.User, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 检查邮箱是否已被其他用户使用
	if user.Email != email {
		existingUser, _ := s.userRepo.GetByEmail(email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, ErrUserAlreadyExists
		}
	}

	// 更新用户信息
	user.Email = email

	err = s.userRepo.Update(user)
	if err != nil {
		logger.Errorf("Failed to update user: %v", err)
		return nil, err
	}

	return user, nil
}

// generateToken 生成JWT token
func (s *userServiceImpl) generateToken(userID uint) (string, error) {
	cfg := config.Load()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(cfg.JWT.ExpiresIn) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.SecretKey))
}

// CreateUserByAdmin 管理员创建用户
func (s *userServiceImpl) CreateUserByAdmin(username, password, email string, status entity.UserStatus) (*entity.User, error) {
	// 唯一性校验
	if u, _ := s.userRepo.GetByUsername(username); u != nil {
		return nil, ErrUserAlreadyExists
	}
	if u, _ := s.userRepo.GetByEmail(email); u != nil {
		return nil, ErrUserAlreadyExists
	}
	// 密码处理
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &entity.User{Username: username, Email: email, PasswordHash: string(hash), Status: status}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser 管理员根据ID获取用户
func (s *userServiceImpl) GetUser(userID uint) (*entity.User, error) {
	return s.GetUserByID(userID)
}

// ListUsers 分页列出用户
func (s *userServiceImpl) ListUsers(page, pageSize int) ([]*entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

// UpdateUserByAdmin 管理员更新用户
func (s *userServiceImpl) UpdateUserByAdmin(userID uint, email string, status entity.UserStatus, password *string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	// 邮箱唯一性
	if email != "" && email != user.Email {
		if exist, _ := s.userRepo.GetByEmail(email); exist != nil && exist.ID != userID {
			return nil, ErrUserAlreadyExists
		}
		user.Email = email
	}
	// 状态更新
	if status != "" {
		user.Status = status
	}
	// 修改密码
	if password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hash)
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser 管理员删除用户
func (s *userServiceImpl) DeleteUser(userID uint) error {
	return s.userRepo.Delete(userID)
}
