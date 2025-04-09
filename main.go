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
	Clipboard string
	Mutex     sync.Mutex
}

var (
	rooms      = make(map[string]*Room)
	roomsMutex sync.Mutex
)

func generateRoomID() string {
	return fmt.Sprintf("%d", time.Now().UnixMicro()/int64(time.Microsecond))
}

func main() {
	r := gin.Default()

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 创建房间
	r.POST("/createRoom", func(c *gin.Context) {
		var req struct {
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Password == "" {
			req.Password = "123456"
		}
		roomID := generateRoomID()
		room := &Room{
			ID:       roomID,
			Password: req.Password,
		}

		roomsMutex.Lock()
		rooms[roomID] = room
		roomsMutex.Unlock()

		c.JSON(http.StatusOK, gin.H{"message": "Room created", "roomID": roomID, "password": room.Password})
	})

	// 加入房间
	r.POST("/joinRoom", func(c *gin.Context) {
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
	r.POST("/clipboard/:roomID", func(c *gin.Context) {
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
		room.Clipboard = content
		room.Mutex.Unlock()

		c.JSON(http.StatusOK, gin.H{"message": "Clipboard content saved!"})
	})

	// 获取剪贴板内容
	r.GET("/clipboard/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		roomsMutex.Lock()
		room, exists := rooms[roomID]
		roomsMutex.Unlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		room.Mutex.Lock()
		content := room.Clipboard
		room.Mutex.Unlock()

		c.JSON(http.StatusOK, gin.H{"content": content})
	})

	r.GET("/", func(c *gin.Context) {
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
