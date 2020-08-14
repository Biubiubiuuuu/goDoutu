package tools

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

const (
	AccessKey  = "VKYaoHZ9no66HILp2XmBMl4RwkvZNLX6F67ek3Qd"
	SecretKey  = "OhvUGMTl45_DykHfhjBMIpd3IKl9g_Qqae2PaXWI"
	Bucketname = "godoutu"
)

// 获取token
func Token() string {
	putPolicy := storage.PutPolicy{
		Scope: Bucketname,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	return putPolicy.UploadToken(mac)
}

// 上传图片
func Upload(imgUrl string) (err error, url string) {
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	upToken := Token()
	uuidKey := uuid.New()
	str := strings.Split(".", imgUrl)
	key := uuidKey.String() + str[len(str)-1]
	if err := formUploader.PutFile(context.Background(), &ret, upToken, key, imgUrl, nil); err != nil {
		return err, ""
	}
	return nil, ret.Key
}
