package api

import (
	"giligili/service"

	"github.com/gin-gonic/gin"
)

// CreateVideo 视频投稿
func CreateVideo(c *gin.Context) {
	videoService := service.CreateVideoService{}
	if err := c.ShouldBind(&videoService); err == nil {
		res := videoService.Create(CurrentUser(c).ID)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// ShowVideo 视频详情接口
func ShowVideo(c *gin.Context) {
	videoService := service.ShowVideoService{}
	c.JSON(200, videoService.Show(c.Param("id")))
}

// ListVideo 视频列表接口
func ListVideo(c *gin.Context) {
	videoService := service.ListVideoService{}
	if err := c.ShouldBind(&videoService); err == nil {
		c.JSON(200, videoService.List())
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// UpdateVideo 更新视频的接口
func UpdateVideo(c *gin.Context) {
	videoService := service.UpdateVideoService{}
	if err := c.ShouldBind(&videoService); err == nil {
		c.JSON(200, videoService.Update(c.Param("id"), CurrentUser(c).ID))
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// DeleteVideo 删除视频
func DeleteVideo(c *gin.Context) {
	videoService := service.DeleteVideoService{}
	c.JSON(200, videoService.Delete(c.Param("id"), CurrentUser(c).ID))
}
