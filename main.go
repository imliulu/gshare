package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Room struct {
	ID        string
	Password  string
	Clipboard []string
	Mutex     sync.Mutex
}

var (
	rooms      = make(map[string]*Room)
	roomsMutex sync.Mutex
)

func generateRoomID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

func main() {
	r := gin.Default()

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// API路由组
	api := r.Group("/api")
	{
		// 创建房间
		api.POST("/rooms", func(c *gin.Context) {
			var req struct {
				ID       string `json:"id"`
				Password string `json:"password"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 如果没有提供 roomID，则生成一个默认的时间戳作为 roomID
			if req.ID == "" {
				req.ID = generateRoomID()
			}

			// 设置默认密码为 "123" 如果用户没有提供密码
			if req.Password == "" {
				req.Password = "123456"
			}

			room := &Room{
				ID:        req.ID,
				Password:  req.Password,
				Clipboard: []string{},
			}

			roomsMutex.Lock()
			rooms[req.ID] = room
			roomsMutex.Unlock()

			c.JSON(http.StatusCreated, gin.H{"message": "Room created", "roomID": req.ID, "password": req.Password})
		})

		// 加入房间
		api.POST("/rooms/join", func(c *gin.Context) {
			var req struct {
				ID       string `json:"id"`
				Password string `json:"password"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			roomsMutex.Lock()
			room, exists := rooms[req.ID]
			roomsMutex.Unlock()

			if !exists || room.Password != req.Password {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid room ID or password"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Joined room", "roomID": room.ID, "password": room.Password})
		})

		// 粘贴文本内容
		api.POST("/rooms/:roomID/clipboard", func(c *gin.Context) {
			roomID := c.Param("roomID")
			roomsMutex.Lock()
			room, exists := rooms[roomID]
			roomsMutex.Unlock()

			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
				return
			}

			var content string
			if err := c.BindJSON(&content); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			room.Mutex.Lock()
			room.Clipboard = append(room.Clipboard, content)
			room.Mutex.Unlock()

			c.JSON(http.StatusOK, gin.H{"message": "Clipboard content saved!"})
		})

		// 获取剪贴板内容
		api.GET("/rooms/:roomID/clipboard", func(c *gin.Context) {
			roomID := c.Param("roomID")
			roomsMutex.Lock()
			room, exists := rooms[roomID]
			roomsMutex.Unlock()

			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
				return
			}

			room.Mutex.Lock()
			contents := room.Clipboard
			room.Mutex.Unlock()

			c.JSON(http.StatusOK, gin.H{"contents": contents})
		})

		// 清空剪贴板内容
		api.DELETE("/rooms/:roomID/clipboard", func(c *gin.Context) {
			roomID := c.Param("roomID")
			roomsMutex.Lock()
			room, exists := rooms[roomID]
			roomsMutex.Unlock()

			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
				return
			}

			room.Mutex.Lock()
			room.Clipboard = []string{}
			room.Mutex.Unlock()

			c.JSON(http.StatusOK, gin.H{"message": "Clipboard content cleared!"})
		})
	}

	// 配置404页面
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/share/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		roomsMutex.Lock()
		room, exists := rooms[roomID]
		roomsMutex.Unlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		c.HTML(http.StatusOK, "share.html", gin.H{"roomID": room.ID, "password": room.Password})
	})

	r.LoadHTMLGlob("templates/*")

	err := r.Run(":8088")
	if err != nil {
		return
	}
}
