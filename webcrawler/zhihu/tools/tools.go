package tools

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

const (
	AccessKey  = "VKYaoHZ9no66HILp2XmBMl4RwkvZNLX6F67ek3Qd"
	SecretKey  = "OhvUGMTl45_DykHfhjBMIpd3IKl9g_Qqae2PaXWI"
	Bucketname = "godoutu"
	MacPath    = "/Users/gaochao/LTWorks/goDoutu/webcrawler/zhihu/basic/image/"
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
	key := uuidKey.String() + "." + str[len(str)-1]
	if err := formUploader.PutFile(context.Background(), &ret, upToken, key, imgUrl, nil); err != nil {
		return err, ""
	}
	return nil, ret.Key
}

// 下载图片到本地
func Download(url string) string {
	str := strings.Split(url, ".")
	fileName := uuid.New().String() + "." + strings.Replace(str[len(str)-1], "?source=1940ef5c", "", -1)
	res, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	time := time.Now().Format("20060102")
	path := MacPath + time
	if !IsExist(path) {
		CreateDir(path)
	}
	path = path + "/" + fileName
	file, err := os.Create(path)
	if err != nil {
		return ""
	}
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	return path
}

// 文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

// 判断文件夹/文件是否存在
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
