package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gotux/config"
	"gotux/middleware"
	"gotux/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取用户信息和配额设置
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 检查存储配额
	storageUsed, err := models.GetUserStorageUsed(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取存储使用量失败"})
		return
	}

	// 解析表单
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "表单解析失败"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有上传文件"})
		return
	}

	// 计算本次上传总大小
	var totalUploadSize int64
	for _, file := range files {
		totalUploadSize += file.Size
	}

	// 检查是否超过配额
	if user.StorageQuota > 0 && storageUsed+totalUploadSize > user.StorageQuota {
		remainingQuota := user.StorageQuota - storageUsed
		c.JSON(http.StatusForbidden, gin.H{
			"error":           "存储空间不足",
			"storage_used":    storageUsed,
			"storage_quota":   user.StorageQuota,
			"remaining_quota": remainingQuota,
			"upload_size":     totalUploadSize,
		})
		return
	}

	var uploadedImages []models.Image
	var errors []string

	for _, file := range files {
		// 检查文件大小
		if file.Size > config.AppConfig.Upload.MaxSize {
			errors = append(errors, fmt.Sprintf("%s: 文件大小超过限制", file.Filename))
			continue
		}

		// 检查文件类型
		if !isAllowedFileType(file.Header.Get("Content-Type")) {
			errors = append(errors, fmt.Sprintf("%s: 不支持的文件类型", file.Filename))
			continue
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: 文件打开失败", file.Filename))
			continue
		}

		// 计算文件哈希
		hash := md5.New()
		if _, err := io.Copy(hash, src); err != nil {
			src.Close()
			errors = append(errors, fmt.Sprintf("%s: 哈希计算失败", file.Filename))
			continue
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))
		src.Close()

		// 检查是否已存在相同文件
		existingImage, err := models.GetImageByHash(hashStr, userID)
		if err == nil {
			uploadedImages = append(uploadedImages, *existingImage)
			continue
		}

		// 生成唯一文件名
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		// 按日期组织文件夹
		dateFolder := time.Now().Format("2006/01/02")
		fullPath := filepath.Join(config.AppConfig.Upload.StoragePath, dateFolder)

		// 创建目录
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			errors = append(errors, fmt.Sprintf("%s: 创建目录失败", file.Filename))
			continue
		}

		// 保存文件
		filePath := filepath.Join(fullPath, newFileName)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			errors = append(errors, fmt.Sprintf("%s: 文件保存失败", file.Filename))
			continue
		}

		// 获取图片尺寸
		width, height := getImageDimensions(filePath)

		// 创建数据库记录
		image := models.Image{
			UserID:       userID,
			FileName:     newFileName,
			OriginalName: file.Filename,
			FilePath:     filepath.Join(dateFolder, newFileName),
			FileSize:     file.Size,
			MimeType:     file.Header.Get("Content-Type"),
			Width:        width,
			Height:       height,
			Hash:         hashStr,
			IsPublic:     true,
		}

		if err := models.CreateImage(&image); err != nil {
			os.Remove(filePath) // 删除已保存的文件
			errors = append(errors, fmt.Sprintf("%s: 数据库保存失败", file.Filename))
			continue
		}

		uploadedImages = append(uploadedImages, image)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("成功上传 %d 个文件", len(uploadedImages)),
		"images":  uploadedImages,
		"errors":  errors,
	})
}

// GetImages 获取图片列表
func GetImages(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var images []models.Image
	var total int64
	var err error

	if keyword != "" {
		images, total, err = models.SearchImages(userID, keyword, page, pageSize)
	} else {
		images, total, err = models.GetImagesByUserID(userID, page, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images":    images,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetImageDetail 获取图片详情
func GetImageDetail(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 检查权限
	if image.UserID != userID {
		user, _ := middleware.GetUser(c)
		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问该图片"})
			return
		}
	}

	// 增加访问次数
	models.IncrementViewCount(uint(imageID))

	c.JSON(http.StatusOK, image)
}

// UpdateImage 更新图片信息
func UpdateImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 检查权限
	if image.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限编辑该图片"})
		return
	}

	var req struct {
		Description string `json:"description"`
		Tags        string `json:"tags"`
		IsPublic    *bool  `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	image.Description = req.Description
	image.Tags = req.Tags
	if req.IsPublic != nil {
		image.IsPublic = *req.IsPublic
	}

	if err := image.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"image":   image,
	})
}

// DeleteImage 删除图片
func DeleteImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 检查权限
	if image.UserID != userID {
		user, _ := middleware.GetUser(c)
		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除该图片"})
			return
		}
	}

	// 删除文件
	fullPath := filepath.Join(config.AppConfig.Upload.StoragePath, image.FilePath)
	os.Remove(fullPath)

	// 删除数据库记录
	if err := image.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// BatchDeleteImages 批量删除图片
func BatchDeleteImages(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req struct {
		ImageIDs []uint `json:"image_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	deletedCount := 0
	for _, imageID := range req.ImageIDs {
		image, err := models.GetImageByID(imageID)
		if err != nil {
			continue
		}

		// 检查权限
		if image.UserID != userID {
			user, _ := middleware.GetUser(c)
			if !user.IsAdmin() {
				continue
			}
		}

		// 删除文件
		fullPath := filepath.Join(config.AppConfig.Upload.StoragePath, image.FilePath)
		os.Remove(fullPath)

		// 删除数据库记录
		if err := image.Delete(); err == nil {
			deletedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("成功删除 %d 张图片", deletedCount),
		"deleted_count": deletedCount,
	})
}

// GetImageLinks 获取图片链接（多种格式）
func GetImageLinks(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	image, err := models.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 获取用户设置
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户设置失败"})
		return
	}

	// 构建完整URL - 使用自定义域名或默认域名
	var baseURL string
	if user.CustomDomain != "" {
		baseURL = user.CustomDomain
	} else {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}
	
	// 使用UUID生成安全链接
	imageURL := fmt.Sprintf("%s/i/%s", baseURL, image.UUID)
	// 兼容旧的直接路径访问
	directURL := fmt.Sprintf("%s/uploads/%s", baseURL, image.FilePath)

	// 生成各种格式的链接
	links := map[string]string{
		"url":                imageURL,
		"direct_url":         directURL,
		"html":               fmt.Sprintf(`<img src="%s" alt="%s" />`, imageURL, image.OriginalName),
		"markdown":           fmt.Sprintf(`![%s](%s)`, image.OriginalName, imageURL),
		"bbcode":             fmt.Sprintf(`[img]%s[/img]`, imageURL),
		"markdown_with_link": fmt.Sprintf(`[![%s](%s)](%s)`, image.OriginalName, imageURL, imageURL),
	}

	c.JSON(http.StatusOK, gin.H{
		"image": image,
		"links": links,
	})
}

// GetImageByUUID 通过UUID获取图片信息
func GetImageByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	
	image, err := models.GetImageByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 增加浏览量
	if image.Stats != nil {
		image.Stats.ViewCount++
		models.DB.Save(image.Stats)
	} else {
		stats := &models.ImageStats{
			ImageID:   image.ID,
			ViewCount: 1,
		}
		models.DB.Create(stats)
	}

	c.JSON(http.StatusOK, gin.H{
		"image": image,
	})
}

// ServeImageByUUID 通过UUID提供图片文件
func ServeImageByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	
	image, err := models.GetImageByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 检查图片是否公开或者是否为图片所有者
	if !image.IsPublic {
		userID, exists := middleware.GetUserID(c)
		if !exists || userID != image.UserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此图片"})
			return
		}
	}

	// 增加浏览量
	if image.Stats != nil {
		image.Stats.ViewCount++
		models.DB.Save(image.Stats)
	} else {
		stats := &models.ImageStats{
			ImageID:   image.ID,
			ViewCount: 1,
		}
		models.DB.Create(stats)
	}

	// 提供文件
	filePath := fmt.Sprintf("uploads/%s", image.FilePath)
	c.File(filePath)
}

// 辅助函数

func isAllowedFileType(mimeType string) bool {
	for _, allowed := range config.AppConfig.Upload.AllowedTypes {
		if allowed == mimeType {
			return true
		}
	}
	return false
}

func getImageDimensions(filePath string) (int, int) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return 0, 0
	}
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// GetStats 获取统计信息
func GetStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取用户信息(包含配额)
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 获取图片总数
	var imageCount int64
	models.DB.Model(&models.Image{}).Where("user_id = ?", userID).Count(&imageCount)

	// 获取存储使用量
	storageUsed, _ := models.GetUserStorageUsed(userID)

	// 获取总访问量
	var totalViews int64
	models.DB.Table("image_stats").
		Joins("JOIN images ON images.id = image_stats.image_id").
		Where("images.user_id = ?", userID).
		Select("COALESCE(SUM(image_stats.view_count), 0)").
		Scan(&totalViews)

	// 计算配额使用百分比和剩余空间
	var quotaPercent float64
	var remainingQuota int64
	if user.StorageQuota > 0 {
		quotaPercent = float64(storageUsed) / float64(user.StorageQuota) * 100
		remainingQuota = user.StorageQuota - storageUsed
	} else {
		quotaPercent = 0
		remainingQuota = -1 // -1 表示无限制
	}

	c.JSON(http.StatusOK, gin.H{
		"image_count":      imageCount,
		"storage_used":     storageUsed,
		"storage_quota":    user.StorageQuota,
		"remaining_quota":  remainingQuota,
		"quota_percent":    quotaPercent,
		"total_views":      totalViews,
	})
}

// GetRandomImage 获取随机图片信息(JSON)
func GetRandomImage(c *gin.Context) {
	query := models.DB.Where("is_public = ?", true)

	// 支持按用户ID筛选
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 支持按标签筛选
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	// 随机获取一张图片
	var image models.Image
	if err := query.Order("RANDOM()").First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到符合条件的图片"})
		return
	}

	// 加载图片统计信息
	models.DB.Model(&image).Association("Stats").Find(&image.Stats)

	c.JSON(http.StatusOK, image)
}

// ServeRandomImage 直接返回随机图片文件(用于图床API)
func ServeRandomImage(c *gin.Context) {
	query := models.DB.Where("is_public = ?", true)

	// 支持按用户ID筛选
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 支持按标签筛选
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	// 随机获取一张图片
	var image models.Image
	if err := query.Order("RANDOM()").First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到符合条件的图片"})
		return
	}

	// 增加浏览次数
	models.IncrementViewCount(image.ID)

	// 构建文件完整路径
	fullPath := filepath.Join(".", image.FilePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片文件不存在"})
		return
	}

	// 设置响应头
	c.Header("Content-Type", image.MimeType)
	c.Header("Cache-Control", "public, max-age=3600")
	c.Header("X-Image-UUID", image.UUID)
	c.Header("X-Image-ID", strconv.Itoa(int(image.ID)))

	// 返回图片文件
	c.File(fullPath)
}

// RedirectRandomImage 重定向到随机图片(用于外部引用)
func RedirectRandomImage(c *gin.Context) {
	query := models.DB.Where("is_public = ?", true)

	// 支持按用户ID筛选
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 支持按标签筛选
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	// 随机获取一张图片
	var image models.Image
	if err := query.Order("RANDOM()").First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到符合条件的图片"})
		return
	}

	// 增加浏览次数
	models.IncrementViewCount(image.ID)

	// 获取用户设置
	var user models.User
	if err := models.DB.First(&user, image.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 构建图片URL
	var baseURL string
	if user.CustomDomain != "" {
		baseURL = user.CustomDomain
	} else {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}

	imageURL := fmt.Sprintf("%s/i/%s", baseURL, image.UUID)

	// 重定向到图片
	c.Redirect(http.StatusFound, imageURL)
}
