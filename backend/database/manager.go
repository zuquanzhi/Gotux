package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseManager struct {
	MainDB *sql.DB
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

// InitMainDB 初始化主数据库（用于用户管理）
func (dm *DatabaseManager) InitMainDB() error {
	db, err := sql.Open("sqlite3", "./gotux.db")
	if err != nil {
		return err
	}
	dm.MainDB = db

	// 创建用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TEXT DEFAULT (datetime('now'))
	);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	return nil
}

// GetUserDB 获取用户专属数据库连接
func (dm *DatabaseManager) GetUserDB(userID int64) (*sql.DB, error) {
	dbPath := fmt.Sprintf("./user_databases/user_%d.db", userID)
	
	// 确保用户数据库目录存在
	if err := os.MkdirAll("./user_databases", 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建用户专属的图片表
	createImagesTable := `
	CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL,
		original_name TEXT NOT NULL,
		size INTEGER NOT NULL,
		mime_type TEXT NOT NULL,
		url TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`

	_, err = db.Exec(createImagesTable)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// CreateUserDatabase 为新用户创建专属数据库
func (dm *DatabaseManager) CreateUserDatabase(userID int64) error {
	db, err := dm.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Printf("Created database for user %d", userID)
	return nil
}

// CloseMainDB 关闭主数据库连接
func (dm *DatabaseManager) CloseMainDB() error {
	if dm.MainDB != nil {
		return dm.MainDB.Close()
	}
	return nil
}