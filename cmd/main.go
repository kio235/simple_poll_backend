package main

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := os.Getenv("MYSQL_DSN") // 从环境变量读取
	// 增加连接重试
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic("failed to connect database after retries")
	}

	// 自动迁移表结构
	db.AutoMigrate(&Question{}, &Option{})

	// 预置数据
	var count int64
	db.Model(&Question{}).Count(&count)
	if count == 0 {
		q := Question{
			Content: "What's your favorite programming language?",
			Options: []Option{
				{Text: "Go"},
				{Text: "JavaScript"},
				{Text: "Python"},
				{Text: "Java"},
			},
		}
		db.Create(&q)
	}
	return db
}

func main() {
	db := initDB()
	broker := NewBroker()

	r := gin.Default()

	// CORS 配置
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// 获取问卷数据
	r.GET("/api/poll", func(c *gin.Context) {
		var question Question
		db.Preload("Options").First(&question)
		c.JSON(200, question)
	})

	// 提交投票
	r.POST("/api/poll/vote/:optionID", func(c *gin.Context) {
		optionID, _ := strconv.Atoi(c.Param("optionID"))

		db.Transaction(func(tx *gorm.DB) error {
			var option Option
			if err := tx.First(&option, optionID).Error; err != nil {
				return err
			}
			return tx.Model(&option).Update("votes", option.Votes+1).Error
		})

		// 广播更新
		var question Question
		db.Preload("Options").First(&question)
		broker.Broadcast(question.Options)
		c.Status(200)
	})

	// SSE 端点
	r.GET("/sse/poll", func(c *gin.Context) {
		c.Stream(func(w io.Writer) bool {
			ch := make(chan string)
			broker.AddClient(ch)
			defer broker.RemoveClient(ch)

			for msg := range ch {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	r.Run(":8080")
}
