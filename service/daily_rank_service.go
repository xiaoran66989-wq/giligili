package service

import (
	"giligili/cache"
	"giligili/model"
	"giligili/serializer"
	"sort"
)

// DailyRankService 每日排行服务
type DailyRankService struct{}

// Get retrieves the most viewed videos for the current process day.
func (service *DailyRankService) Get() serializer.Response {
	ids := cache.DailyRankIDs(10)
	videos := []model.Video{}
	if len(ids) > 0 {
		if err := model.DB.Where("id IN (?)", ids).Find(&videos).Error; err != nil {
			return serializer.Response{Status: 50000, Msg: "数据库连接错误", Error: err.Error()}
		}
		positions := make(map[uint]int, len(ids))
		for index, id := range ids {
			positions[id] = index
		}
		sort.Slice(videos, func(i, j int) bool {
			return positions[videos[i].ID] < positions[videos[j].ID]
		})
	}
	return serializer.Response{Data: serializer.BuildVideos(videos)}
}
