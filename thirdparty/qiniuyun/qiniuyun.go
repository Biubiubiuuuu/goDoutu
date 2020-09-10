package qiniuyun

import (
	"context"
	"strings"

	"github.com/Biubiubiuuuu/goDoutu/helper/config"
	"github.com/google/uuid"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

// 七牛云获取token
func Token() string {
	putPolicy := storage.PutPolicy{
		Scope: config.QNYBucketname,
	}
	mac := qbox.NewMac(config.QNYAccessKey, config.QNYSecretKey)
	return putPolicy.UploadToken(mac)
}

// 七牛云上传图片
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
	str := strings.Split(imgUrl, ".")
	key := uuidKey.String() + "." + str[len(str)-1]
	if err := formUploader.PutFile(context.Background(), &ret, upToken, key, imgUrl, nil); err != nil {
		return err, ""
	}
	return nil, config.QNYBasicUrl + ret.Key
}
