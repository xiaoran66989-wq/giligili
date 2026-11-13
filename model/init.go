package model

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB 数据库链接单例
var DB *gorm.DB

// Database 在中间件中初始化mysql链接
func Database() {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		if runtime.GOOS == "windows" {
			driver = "sqlite3"
		} else {
			driver = "mysql"
		}
	}

	connString := os.Getenv("MYSQL_DSN")
	if driver == "sqlite3" {
		connString = os.Getenv("SQLITE_PATH")
		if connString == "" {
			configDir, err := os.UserConfigDir()
			if err != nil {
				panic(err)
			}
			dataDir := filepath.Join(configDir, "Giligili")
			if err := os.MkdirAll(dataDir, 0o755); err != nil {
				panic(err)
			}
			connString = filepath.Join(dataDir, "giligili.db")
		}
		connString = "file:" + filepath.ToSlash(connString) + "?_busy_timeout=5000&_journal_mode=WAL"
	}

	db, err := gorm.Open(driver, connString)
	if err != nil {
		panic(err)
	}

	db.LogMode(gin.Mode() != gin.ReleaseMode)
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(20)
	if driver == "sqlite3" {
		sqlDB.SetMaxOpenConns(1)
	} else {
		sqlDB.SetMaxOpenConns(100)
	}
	sqlDB.SetConnMaxLifetime(30 * time.Second)

	DB = db
	migration(driver)
}
