package osssdk

import (
	"bytes"
	"context"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
	"net/url"
)

var _ OSSModel = (*AliYunOssModel)(nil)

type OssConf struct {
	OssAccessKeyId     string
	OssAccessKeySecret string
	OssEndpoint        string
	BucketName         string
	IsCname            bool
}
type OSSModel interface {
	UploadOSSByUrl(ctx context.Context, objectKey string, url string) error
	UploadOssByFile(ctx context.Context, objectKey string, data []byte) error
	ParseUrl(ctx context.Context, url string) (string, error)
	GetSignUrl(ctx context.Context, objectKey string, expiredInSec int64, options ...oss.Option) (string, error)
}

type AliYunOssModel struct {
	Bucket string
	Client *oss.Client
}

func NewAliYunOssModel(ossConf OssConf) OSSModel {
	accessKeyId := ossConf.OssAccessKeyId
	accessKeySecret := ossConf.OssAccessKeySecret
	endpoint := ossConf.OssEndpoint
	// 创建OSS客户端实例
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret, oss.UseCname(ossConf.IsCname))
	if err != nil {
		return nil
	}
	aliYunOssModel := &AliYunOssModel{
		Bucket: ossConf.BucketName,
		Client: client,
	}
	return aliYunOssModel
}

func (o *AliYunOssModel) UploadOSSByUrl(ctx context.Context, objectKey string, urlString string) error {
	response, err := http.Get(urlString)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	bucket, err := o.Client.Bucket(o.Bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObject(objectKey, response.Body)
	if err != nil {
		return err
	}
	return nil
}
func (o *AliYunOssModel) UploadOssByFile(ctx context.Context, objectKey string, data []byte) error {
	if len(data) <= 0 {
		return errors.New("数据不能为空")
	}
	bucket, err := o.Client.Bucket(o.Bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObject(objectKey, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}
func (o *AliYunOssModel) ParseUrl(ctx context.Context, fileUrl string) (string, error) {
	parsedURL, err := url.Parse(fileUrl)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}
func (o *AliYunOssModel) GetSignUrl(ctx context.Context, objectKey string, expiredInSec int64, options ...oss.Option) (string, error) {
	bucket, err := o.Client.Bucket(o.Bucket)
	if err != nil {
		return "", err
	}
	signedURL, err := bucket.SignURL(objectKey, oss.HTTPGet, expiredInSec, options...)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}
