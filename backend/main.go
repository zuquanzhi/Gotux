// backend/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// ========== 数据模型 ==========
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Image struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	URL          string `json:"url"`
	CreatedAt    string `json:"created_at"`
}

// ========== JWT 配置 ==========
var jwtSecret = []byte("gotux-secret-key-for-development")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ========== 工具函数 ==========
func generateJWT(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

// ========== 数据库初始化 ==========
func initDB() *sql.DB {
	// 创建 uploads 目录
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", "./gotux.db")
	if err != nil {
		log.Fatal(err)
	}

	// 创建用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TEXT DEFAULT (datetime('now'))
	);`

	// 创建图片表
	createImagesTable := `
	CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER DEFAULT 0,
		filename TEXT NOT NULL,
		original_name TEXT NOT NULL,
		size INTEGER NOT NULL,
		mime_type TEXT NOT NULL,
		url TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createImagesTable)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// ========== 处理函数 ==========
func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 检查用户是否已存在
		var existingUser User
		err := db.QueryRow("SELECT id FROM users WHERE username = ? OR email = ?",
			user.Username, user.Email).Scan(&existingUser.ID)

		if err == nil {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}

		// 创建新用户
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		result, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			user.Username, user.Email, hashedPassword)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		userID, _ := result.LastInsertId()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"user": map[string]interface{}{
				"id":       userID,
				"username": user.Username,
				"email":    user.Email,
			},
		})
	}
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 获取用户
		var user User
		err := db.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?",
			credentials.Username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// 验证密码
		if !checkPasswordHash(credentials.Password, user.Password) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// 生成 JWT token
		token, err := generateJWT(user.ID, user.Username)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"token":   token,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		})
	}
}

func uploadImageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析表单
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

		// 验证文件类型
		if !isValidImageType(header.Header.Get("Content-Type")) {
			http.Error(w, "Invalid file type. Only images allowed", http.StatusBadRequest)
			return
		}

		// 生成唯一文件名
		timestamp := time.Now().UnixNano()
		extension := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d%s", timestamp, extension)

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

		// 创建图片记录（公开上传，user_id = 0）
		image := Image{
			UserID:       0,
			Filename:     filename,
			OriginalName: header.Filename,
			Size:         header.Size,
			MimeType:     header.Header.Get("Content-Type"),
			URL:          "http://localhost:8080/uploads/" + filename,
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}

		// 保存到数据库
		query := `INSERT INTO images (user_id, filename, original_name, size, mime_type, url, created_at) 
				  VALUES (?, ?, ?, ?, ?, ?, ?)`
		result, err := db.Exec(query, image.UserID, image.Filename, image.OriginalName,
			image.Size, image.MimeType, image.URL, image.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to save image info", http.StatusInternalServerError)
			return
		}

		image.ID, _ = result.LastInsertId()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    image,
		})
	}
}

func getMyImagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 简化：为了演示，我们假设用户ID为1
		// 实际项目中应该从JWT token解析用户ID
		userID := int64(1)

		rows, err := db.Query("SELECT id, user_id, filename, original_name, size, mime_type, url, created_at FROM images WHERE user_id = ? ORDER BY created_at DESC", userID)
		if err != nil {
			http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var images []Image
		for rows.Next() {
			var img Image
			err := rows.Scan(&img.ID, &img.UserID, &img.Filename, &img.OriginalName,
				&img.Size, &img.MimeType, &img.URL, &img.CreatedAt)
			if err != nil {
				http.Error(w, "Failed to scan image", http.StatusInternalServerError)
				return
			}
			images = append(images, img)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    images,
		})
	}
}

func main() {
	// 初始化数据库
	db := initDB()
	defer db.Close()

	// 设置路由
	http.HandleFunc("/api/auth/register", registerHandler(db))
	http.HandleFunc("/api/auth/login", loginHandler(db))
	http.HandleFunc("/api/images/upload", uploadImageHandler(db))
	http.HandleFunc("/api/images/my", getMyImagesHandler(db))

	// 静态文件服务（图片）
	http.HandleFunc("/uploads/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "."+r.URL.Path)
	})

	// 前端页面路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "../frontend/index.html")
		} else if r.URL.Path == "/dashboard" {
			http.ServeFile(w, r, "../frontend/dashboard.html")
		} else {
			http.NotFound(w, r)
		}
	})

	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("📁 Upload page: http://localhost:8080")
	fmt.Println("📊 Dashboard: http://localhost:8080/dashboard")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
