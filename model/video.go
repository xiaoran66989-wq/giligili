package model

import (
	"giligili/cache"
	"net/url"
	"os"

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
	Views  uint64
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

// View returns the persisted total view count.
func (video *Video) View() uint64 {
	return video.Views
}

// AddView increments the persisted total and updates the in-process daily rank.
func (video *Video) AddView() {
	if err := DB.Model(video).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error; err != nil {
		return
	}
	video.Views++
	cache.RecordDailyView(video.ID)
}
