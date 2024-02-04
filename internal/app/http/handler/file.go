package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"social/config"
	"social/ent"
	"social/internal/usecase"
	"social/pkg/files"
	"social/pkg/log"
)

type FileHandler struct {
	userUsecase *usecase.UserUseCase
}

func NewFileHandler(userUsecase *usecase.UserUseCase) *FileHandler {
	return &FileHandler{
		userUsecase: userUsecase,
	}
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

func (h *FileHandler) UpdateAvatar(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	base64Data := c.PostForm("image")
	path, err := files.ImgFromBase64("", base64Data)
	if err != nil {
		resp.Code = 1
		resp.Message = "user not exists"
		c.JSON(http.StatusOK, resp)
		return
	}
	_, err = h.userUsecase.UpdateFieldUser(c, u.ID, path)
	if err != nil {
		resp.Code = 1
		resp.Message = "update not image fail"
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "ok"
	c.JSON(http.StatusOK, resp)
}
