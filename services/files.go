package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadFile(c *gin.Context) {
	roomID := c.Param("roomID")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上传失败：" + err.Error()})
		return
	}

	// 生成新文件名：时间戳 + 原扩展名
	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := fmt.Sprintf("static/uploads/%s/%s", roomID, newFilename)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	// 保存文件
	dst, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件写入失败"})
		return
	}

	// 返回文件信息
	c.JSON(http.StatusOK, gin.H{
		"url":     "/uploads/" + roomID + "/" + newFilename,
		"name":    header.Filename,
		"size":    header.Size,
		"type":    header.Header.Get("Content-Type"),
		"isImage": isImage(header.Header.Get("Content-Type")),
	})
}

// 判断是否为图片类型
func isImage(mimeType string) bool {
	imageTypes := []string{"image/jpeg", "image/png", "image/gif"}
	for _, t := range imageTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}

func ListFiles(c *gin.Context) {
	roomID := c.Param("roomID")
	dirPath := fmt.Sprintf("static/uploads/%s", roomID)

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"files": []gin.H{}})
		return
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "读取目录失败"})
		return
	}

	var fileList []gin.H
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		info, err := entry.Info()
		if err != nil {
			// 处理获取文件信息失败的情况（如跳过该文件）
			continue
		}
		fileList = append(fileList, gin.H{
			"name":    filename,
			"url":     fmt.Sprintf("/static/uploads/%s/%s", roomID, filename),
			"size":    info.Size(), // 现在是安全的
			"isImage": true,
		})
	}

	c.JSON(http.StatusOK, gin.H{"files": fileList})
}

// 根据文件扩展名判断是否为图片
//func isImageByExt(ext string) bool {
//	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
//	for _, e := range imageExtensions {
//		if ext == e {
//			return true
//		}
//	}
//	return false
