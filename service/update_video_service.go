package service

import (
	"giligili/model"
	"giligili/serializer"
)

// UpdateVideoService 更新视频的服务
type UpdateVideoService struct {
	Title string `form:"title" json:"title" binding:"required,min=2,max=30"`
	Info  string `form:"info" json:"info" binding:"max=300"`
}

// Update 更新视频
func (service *UpdateVideoService) Update(id string, userID uint) serializer.Response {
	var video model.Video
	if err := model.DB.Where("id = ? AND user_id = ?", id, userID).First(&video).Error; err != nil {
		return serializer.Response{Status: 404, Msg: "视频不存在或无权修改"}
	}

	video.Title = service.Title
	video.Info = service.Info
	if err := model.DB.Save(&video).Error; err != nil {
		return serializer.Response{Status: 50002, Msg: "视频保存失败"}
	}

	return serializer.Response{Data: serializer.BuildVideo(video)}
}
