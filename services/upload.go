package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadFiles(c *gin.Context) {
	roomID := c.Param("roomID")

	// 解析 multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析表单失败: " + err.Error()})
		return
	}

	// 获取所有上传的文件
	files := c.Request.MultipartForm.File["file"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有文件上传"})
		return
	}

	var uploadedFiles []gin.H

	for _, file := range files {
		// 生成新文件名：时间戳 + 原扩展名
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := fmt.Sprintf("static/uploads/%s/%s", roomID, newFilename)

		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
			return
		}

		// 保存文件
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败: " + err.Error()})
			return
		}
		defer src.Close()

		dst, err := os.Create(savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败: " + err.Error()})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件写入失败: " + err.Error()})
			return
		}

		// 返回文件信息
		uploadedFiles = append(uploadedFiles, gin.H{
			"url":     "/uploads/" + roomID + "/" + newFilename,
			"name":    file.Filename,
			"size":    file.Size,
			"type":    file.Header.Get("Content-Type"),
			"isImage": isImageByExt(ext),
		})
	}

	c.JSON(http.StatusOK, gin.H{"files": uploadedFiles})
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
		ext := filepath.Ext(filename)
		info, err := entry.Info()
		if err != nil {
			// 处理获取文件信息失败的情况（如跳过该文件）
			continue
		}
		fileList = append(fileList, gin.H{
			"name":    filename,
			"url":     fmt.Sprintf("/static/uploads/%s/%s", roomID, filename),
			"size":    info.Size(),
			"isImage": isImageByExt(ext),
		})
	}

	c.JSON(http.StatusOK, gin.H{"files": fileList})
}

// 根据文件扩展名判断是否为图片
func isImageByExt(ext string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	ext = strings.ToLower(ext) // 确保扩展名是小写的
	for _, e := range imageExtensions {
		if ext == e {
			return true
		}
	}
	return false
}
