package service

import (
	"giligili/model"
	"giligili/serializer"
)

// DeleteVideoService 删除投稿的服务
type DeleteVideoService struct{}

// Delete 删除视频
func (service *DeleteVideoService) Delete(id string, userID uint) serializer.Response {
	var video model.Video
	if err := model.DB.Where("id = ? AND user_id = ?", id, userID).First(&video).Error; err != nil {
		return serializer.Response{Status: 404, Msg: "视频不存在或无权删除"}
	}
	if err := model.DB.Delete(&video).Error; err != nil {
		return serializer.Response{Status: 50000, Msg: "视频删除失败"}
	}
	return serializer.Response{}
}
