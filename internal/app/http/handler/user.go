package handler

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"social/config"
	"social/ent"
	"social/internal/app/http/middleware"
	"social/internal/entity"
	"social/internal/usecase"
	"social/pkg/files"
	"social/pkg/log"
)

type UserHandler struct {
	userUsecase *usecase.UserUseCase
}

func NewUserHandler(userUsecase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *UserHandler) Login(c *gin.Context) {
	req := entity.UserLoginReq{}
	resp := new(ApiResponse)
	if err := c.ShouldBind(&req); err != nil {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user, err := h.userUsecase.UserLogin(ctx, req)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	at, rt, err := middleware.GenerateToken(user)
	if err != nil {
		resp.Code = 1
		resp.Message = "generate token fail"
		c.JSON(http.StatusOK, resp)
		return
	}
	result := map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"nickname":      user.Nickname,
		"sex":           user.Sex,
		"avatar":        user.Avatar,
		"access_token":  at,
		"refresh_token": rt,
	}
	resp.Message = "login successfully"
	resp.Data = result
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Register(c *gin.Context) {
	req := entity.UserRegisterReq{}
	resp := new(ApiResponse)
	if err := c.ShouldBind(&req); err != nil {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		return
	}
	m := make(map[string]interface{})
	if err := copier.Copy(&m, req); err != nil {
		fmt.Println(err)
		resp.Code = 1
		resp.Message = "params invalid."
		c.JSON(http.StatusOK, resp)
		return
	}
	//filter.Filters(m)
	if err := h.userUsecase.UserRegister(c, req); err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Message = "register successfully"
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Info(c *gin.Context) {
	resp := new(ApiResponse)
	// Get userinfo from token
	user, exists := c.Get("jwt-user")
	if !exists {
		resp.Code = 1
		resp.Message = "user not exists"
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Message = "ok"
	resp.Data = user
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) RetrievePassword(c *gin.Context) {

}

func (h *UserHandler) GetEmoji(c *gin.Context) {

}

func (h *UserHandler) UploadFile(c *gin.Context) {
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

func (h *UserHandler) UpdateAvatar(c *gin.Context) {
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
	_, err = h.userUsecase.UpdateFieldUserRepo(c, u.ID, path)
	if err != nil {
		resp.Code = 1
		resp.Message = "update not image fail"
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "ok"
	c.JSON(http.StatusOK, resp)
}
