package handlers

import (
	"encoding/json"
	"fmt"
	"gotux/database"
	"gotux/models"
	"gotux/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageHandler struct {
	DBManager *database.DatabaseManager
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("UploadImage called: Method=%s, URL=%s\n", r.Method, r.URL.Path)
	
	if r.Method != "POST" {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从JWT获取用户ID（支持匿名上传）
	fmt.Printf("Getting authorization header...\n")
	var userID int64 = 1 // 默认用户ID为1（匿名用户）
	
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		fmt.Printf("Parsing JWT token...\n")
		tokenString := authHeader[7:]
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			fmt.Printf("JWT parsing failed: %v, using anonymous user\n", err)
		} else {
			userID = claims.UserID
			fmt.Printf("User ID extracted: %d\n", userID)
		}
	} else {
		fmt.Printf("No authorization header, using anonymous user (ID: %d)\n", userID)
	}

	err := r.ParseMultipartForm(50 << 20) // 增加到50MB以支持多文件
	if err != nil {
		fmt.Printf("Failed to parse multipart form: %v\n", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// 获取用户专属数据库
	userDB, err := h.DBManager.GetUserDB(userID)
	if err != nil {
		fmt.Printf("Error getting user database for user %d: %v\n", userID, err)
		http.Error(w, "Failed to access user database", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	// 创建用户专属的上传目录
	userUploadDir := fmt.Sprintf("./uploads/user_%d", userID)
	if _, err := os.Stat(userUploadDir); os.IsNotExist(err) {
		os.MkdirAll(userUploadDir, 0755)
	}

	// 处理单个或多个文件
	files := r.MultipartForm.File["image"]
	if len(files) == 0 {
		fmt.Printf("No files uploaded\n")
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}

	var uploadedImages []models.Image
	var errors []string

	for _, header := range files {
		fmt.Printf("Processing file: %s, Size: %d\n", header.Filename, header.Size)

		file, err := header.Open()
		if err != nil {
			fmt.Printf("Failed to open file %s: %v\n", header.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to open file %s", header.Filename))
			continue
		}

		if !isValidImageType(header.Header.Get("Content-Type")) {
			file.Close()
			errors = append(errors, fmt.Sprintf("Invalid file type for %s", header.Filename))
			continue
		}

		// 生成唯一文件名
		timestamp := time.Now().UnixNano()
		extension := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d%s", timestamp, extension)

		// 保存文件到用户专属目录
		filePath := filepath.Join(userUploadDir, filename)
		out, err := os.Create(filePath)
		if err != nil {
			file.Close()
			errors = append(errors, fmt.Sprintf("Failed to save file %s", header.Filename))
			continue
		}

		_, err = io.Copy(out, file)
		out.Close()
		file.Close()

		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to save file %s", header.Filename))
			continue
		}

		// 保存到用户专属数据库
		image := models.Image{
			Filename:     filename,
			OriginalName: header.Filename,
			Size:         header.Size,
			MimeType:     header.Header.Get("Content-Type"),
			URL:          fmt.Sprintf("http://localhost:8081/uploads/user_%d/%s", userID, filename),
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}

		query := `INSERT INTO images (filename, original_name, size, mime_type, url, created_at) 
				  VALUES (?, ?, ?, ?, ?, ?)`
		result, err := userDB.Exec(query, image.Filename, image.OriginalName,
			image.Size, image.MimeType, image.URL, image.CreatedAt)
		if err != nil {
			fmt.Printf("Error inserting image into database for user %d: %v\n", userID, err)
			errors = append(errors, fmt.Sprintf("Failed to save image info for %s", header.Filename))
			continue
		}

		image.ID, _ = result.LastInsertId()
		uploadedImages = append(uploadedImages, image)
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": len(uploadedImages) > 0,
		"data":    uploadedImages,
		"total":   len(files),
		"uploaded": len(uploadedImages),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ImageHandler) GetMyImages(w http.ResponseWriter, r *http.Request) {
	// 从JWT获取用户ID
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[7:]
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	// 获取用户专属数据库
	userDB, err := h.DBManager.GetUserDB(userID)
	if err != nil {
		http.Error(w, "Failed to access user database", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	rows, err := userDB.Query("SELECT id, filename, original_name, size, mime_type, url, created_at FROM images ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		err := rows.Scan(&img.ID, &img.Filename, &img.OriginalName,
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
