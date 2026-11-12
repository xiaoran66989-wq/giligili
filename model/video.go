package model

import (
	"giligili/cache"
	"net/url"
	"os"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/jinzhu/gorm"
)

// Video 视频模型
type Video struct {
	gorm.Model
	UserID uint
	Title  string
	Info   string
	URL    string
	Avatar string
}

// mediaURL returns direct HTTP(S) links unchanged and signs OSS object keys
// only when OSS is configured.
func mediaURL(value string) string {
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != "" {
		return value
	}

	endpoint := os.Getenv("OSS_END_POINT")
	accessKeyID := os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")
	bucketName := os.Getenv("OSS_BUCKET")
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return ""
	}

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return ""
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return ""
	}
	signedURL, err := bucket.SignURL(value, oss.HTTPGet, 600)
	if err != nil {
		return ""
	}
	return signedURL
}

// AvatarURL 封面地址
func (video *Video) AvatarURL() string {
	return mediaURL(video.Avatar)
}

// VideoURL 视频地址
func (video *Video) VideoURL() string {
	return mediaURL(video.URL)
}

// View 点击数
func (video *Video) View() uint64 {
	countStr, _ := cache.RedisClient.Get(cache.VideoViewKey(video.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 视频游览
func (video *Video) AddView() {
	cache.RedisClient.Incr(cache.VideoViewKey(video.ID))
	cache.RedisClient.ZIncrBy(cache.DailyRankKey, 1, strconv.Itoa(int(video.ID)))
}
