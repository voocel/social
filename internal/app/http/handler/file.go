package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"social/config"
	"social/ent"
	"social/pkg/files"
	"social/pkg/log"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	_, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		log.Error(err)
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := files.GenFilename(ext)
	folderPath := filepath.Join(config.Conf.App.StaticRootPath, "images", filename)
	if err := c.SaveUploadedFile(file, folderPath); err != nil {
		resp.Code = 1
		resp.Message = "upload file fail"
		c.JSON(http.StatusOK, resp)
		return
	}

	data := map[string]interface{}{
		"url": config.Conf.App.Domain + "/" + folderPath,
	}
	resp.Message = "ok"
	resp.Data = data
	c.JSON(http.StatusOK, resp)
}
