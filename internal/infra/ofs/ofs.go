package ofs

import (
	"context"
	"fmt"
	"gin-starter/config"
	"gin-starter/pkg/utils"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Bucket = struct {
	Demo string
}{
	Demo: "demo",
}

var client *s3.Client

func InitOfs() {
	if config.AppConfig == nil {
		utils.Log.Fatal("配置未初始化")
	}
	accessKeyID := config.AppConfig.S3.AccessKeyID
	secretAccessKey := config.AppConfig.S3.SecretAccessKey
	region := config.AppConfig.S3.Region
	endpoint := config.AppConfig.S3.Endpoint

	if accessKeyID == "" || secretAccessKey == "" || region == "" || endpoint == "" {
		utils.Log.Fatalf("S3初始化失败，参数缺失")
	}
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")
	cfg := aws.Config{
		Region:      region,
		Credentials: creds,
	}

	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		utils.Log.Fatalf("S3连接失败: %v", err)
	}
	bucketList := make([]string, len(resp.Buckets))
	for i, b := range resp.Buckets {
		bucketList[i] = *b.Name
	}
	utils.Log.Info(fmt.Sprintf("Buckets: %s", strings.Join(bucketList, ", ")))
	utils.Log.Info("S3初始化完成")
}

func GetClient() *s3.Client {
	return client
}
