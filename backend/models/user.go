package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"default:'user'" json:"role"` // admin, user
	Avatar    string         `json:"avatar"`
	Status    string         `gorm:"default:'active'" json:"status"` // active, disabled
	Images    []Image        `gorm:"foreignKey:UserID" json:"images,omitempty"`
}

// HashPassword 加密密码
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// CreateUser 创建用户
func CreateUser(username, email, password string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		Role:     "user",
		Status:   "active",
	}

	if err := user.HashPassword(password); err != nil {
		return nil, err
	}

	if err := DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id uint) (*User, error) {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (u *User) Update() error {
	return DB.Save(u).Error
}

// DeleteUser 删除用户
func (u *User) Delete() error {
	return DB.Delete(u).Error
}

// GetAllUsers 获取所有用户
func GetAllUsers(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	if err := DB.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CreateDefaultAdmin 创建默认管理员
func CreateDefaultAdmin() {
	var count int64
	DB.Model(&User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		admin := &User{
			Username: "admin",
			Email:    "admin@gotux.com",
			Role:     "admin",
			Status:   "active",
		}

		if err := admin.HashPassword("admin123"); err != nil {
			return
		}

		if err := DB.Create(admin).Error; err != nil {
			return
		}
	}
}

// IsAdmin 检查是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// ValidateUser 验证用户登录
func ValidateUser(username, password string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !user.IsActive() {
		return nil, errors.New("账户已被禁用")
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}
