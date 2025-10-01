package controllers

import (
	"gotux/config"
	"gotux/middleware"
	"gotux/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 检查用户名是否已存在
	if _, err := models.GetUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if _, err := models.GetUserByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 创建用户
	user, err := models.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user":    user,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证用户
	user, err := models.ValidateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成 JWT Token
	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile 获取当前用户信息
func GetProfile(c *gin.Context) {
	user, exists := middleware.GetUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取存储使用情况
	storageUsed, _ := models.GetUserStorageUsed(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"storage_used": storageUsed,
	})
}

// UpdateProfile 更新用户信息
func UpdateProfile(c *gin.Context) {
	user, exists := middleware.GetUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req struct {
		Email  string `json:"email" binding:"omitempty,email"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user":    user,
	})
}

// UpdateSettings 更新用户设置
func UpdateSettings(c *gin.Context) {
	user, exists := middleware.GetUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req struct {
		CustomDomain      string `json:"custom_domain"`
		DefaultLinkFormat string `json:"default_link_format"`
		EnableWatermark   *bool  `json:"enable_watermark"`
		WatermarkText     string `json:"watermark_text"`
		WatermarkPosition string `json:"watermark_position"`
		CompressImage     *bool  `json:"compress_image"`
		CompressQuality   *int   `json:"compress_quality"`
		MaxImageSize      *int64 `json:"max_image_size"`
		AllowedImageTypes string `json:"allowed_image_types"`
		EnableImageReview *bool  `json:"enable_image_review"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 更新设置
	if req.CustomDomain != "" {
		user.CustomDomain = req.CustomDomain
	}
	if req.DefaultLinkFormat != "" {
		user.DefaultLinkFormat = req.DefaultLinkFormat
	}
	if req.EnableWatermark != nil {
		user.EnableWatermark = *req.EnableWatermark
	}
	if req.WatermarkText != "" {
		user.WatermarkText = req.WatermarkText
	}
	if req.WatermarkPosition != "" {
		user.WatermarkPosition = req.WatermarkPosition
	}
	if req.CompressImage != nil {
		user.CompressImage = *req.CompressImage
	}
	if req.CompressQuality != nil {
		if *req.CompressQuality < 1 || *req.CompressQuality > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "压缩质量必须在 1-100 之间"})
			return
		}
		user.CompressQuality = *req.CompressQuality
	}
	if req.MaxImageSize != nil {
		if *req.MaxImageSize < 0 || *req.MaxImageSize > 52428800 { // 最大 50MB
			c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小限制必须在 0-50MB 之间"})
			return
		}
		user.MaxImageSize = *req.MaxImageSize
	}
	if req.AllowedImageTypes != "" {
		user.AllowedImageTypes = req.AllowedImageTypes
	}
	if req.EnableImageReview != nil {
		user.EnableImageReview = *req.EnableImageReview
	}

	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "设置更新成功",
		"user":    user,
	})
}

// GetSettings 获取用户设置
func GetSettings(c *gin.Context) {
	user, exists := middleware.GetUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"settings": gin.H{
			"custom_domain":       user.CustomDomain,
			"default_link_format": user.DefaultLinkFormat,
			"enable_watermark":    user.EnableWatermark,
			"watermark_text":      user.WatermarkText,
			"watermark_position":  user.WatermarkPosition,
			"compress_image":      user.CompressImage,
			"compress_quality":    user.CompressQuality,
			"max_image_size":      user.MaxImageSize,
			"allowed_image_types": user.AllowedImageTypes,
			"enable_image_review": user.EnableImageReview,
			"storage_quota":       user.StorageQuota,
			"used_storage":        user.UsedStorage,
		},
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	user, exists := middleware.GetUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证旧密码
	if err := user.CheckPassword(req.OldPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "原密码错误"})
		return
	}

	// 设置新密码
	if err := user.HashPassword(req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// generateToken 生成 JWT Token
func generateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireTime) * time.Hour)

	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}
