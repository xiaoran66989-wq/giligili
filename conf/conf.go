package conf

import (
	"giligili/model"
	"giligili/tasks"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	// Load local settings when present; container environment variables take precedence.
	_ = godotenv.Load()
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	if err := LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
		panic(err)
	}

	model.Database()
	tasks.CronJob()
}
