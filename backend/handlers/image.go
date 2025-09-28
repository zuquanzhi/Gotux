package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gotux/models"
)

type ImageHandler struct {
	DB *sql.DB
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// 公开上传，不需要认证
	if r.Method != "POST" {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !isValidImageType(header.Header.Get("Content-Type")) {
		http.Error(w, "Invalid file type. Only images allowed", http.StatusBadRequest)
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().UnixNano()
	extension := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", timestamp, extension)

	// 确保 uploads 目录存在
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	// 保存文件
	out, err := os.Create("./uploads/" + filename)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = out.ReadFrom(file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// 获取用户ID（如果已登录）
	var userID int64 = 0
	if r.Header.Get("Authorization") != "" {
		// 这里可以解析 JWT 获取用户ID，简化处理
		// 实际项目中应该实现 JWT 验证
	}

	// 保存到数据库
	image := models.Image{
		UserID:       userID,
		Filename:     filename,
		OriginalName: header.Filename,
		Size:         header.Size,
		MimeType:     header.Header.Get("Content-Type"),
		URL:          "http://localhost:8080/uploads/" + filename,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	query := `INSERT INTO images (user_id, filename, original_name, size, mime_type, url, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := h.DB.Exec(query, image.UserID, image.Filename, image.OriginalName,
		image.Size, image.MimeType, image.URL, image.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to save image info", http.StatusInternalServerError)
		return
	}

	image.ID, _ = result.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    image,
	})
}

func (h *ImageHandler) GetMyImages(w http.ResponseWriter, r *http.Request) {
	// 需要认证的路由
	var userID int64 = 1 // 简化：假设用户ID为1，实际应该从JWT解析
	
	// 从Authorization头获取token并解析用户ID
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// 移除 "Bearer " 前缀
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			_ = authHeader[7:] // tokenString 暂时不使用
		}
		
		// 这里应该解析JWT获取真实的用户ID
		// 为了简化，暂时使用固定值
	}

	rows, err := h.DB.Query("SELECT id, user_id, filename, original_name, size, mime_type, url, created_at FROM images WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		err := rows.Scan(&img.ID, &img.UserID, &img.Filename, &img.OriginalName,
			&img.Size, &img.MimeType, &img.URL, &img.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to scan image", http.StatusInternalServerError)
			return
		}
		images = append(images, img)
	}

	// 确保返回空数组而不是null
	if images == nil {
		images = []models.Image{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    images,
	})
}

func isValidImageType(mimeType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[mimeType]
}
