package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pidanou/c1-core/pkg/connector"
	"github.com/pidanou/c1-core/pkg/connector/proto"
)

type S3Connector struct {
	logger   hclog.Logger
	S3Client *s3.Client
}

type Options struct {
	Profile string   `json:"profile"`
	Buckets []string `json:"buckets"`
	Region  string   `json:"region"`
}

func (s *S3Connector) Sync(options string, cb connector.CallbackInterface) (string, error) {

	time.Sleep(time.Second * 20)
	var opts Options

	err := json.Unmarshal([]byte(options), &opts)
	if err != nil {
		s.logger.Info("Failed to unmarshal options", "error", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(opts.Region),
		config.WithSharedConfigProfile(opts.Profile),
	)

	// Create S3 service client
	svc := s3.NewFromConfig(cfg)
	s.S3Client = svc

	var buckets []string
	if opts.Buckets != nil {
		buckets = opts.Buckets
	} else {
		buckets, err = s.listBuckets()
		if err != nil {
			s.logger.Warn("Failed to list buckets", err)
			return fmt.Sprintf("{\"message\":\"%v\"}", err), err
		}
	}

	for _, bucket := range buckets {
		s.listObjects(bucket, opts, cb)
	}
	return `{"message":"success"}`, nil
}

func (s *S3Connector) listBuckets() ([]string, error) {
	res := []string{}
	result, err := s.S3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	for _, bucket := range result.Buckets {
		var noname = ""
		if bucket.Name == nil {
			bucket.Name = &noname
		}
		res = append(res, *bucket.Name)
	}
	return res, nil
}

func (s *S3Connector) listObjects(bucket string, opts Options, cb connector.CallbackInterface) {
	params := &s3.ListObjectsV2Input{
		Bucket: &bucket,
	}
	p := s3.NewListObjectsV2Paginator(s.S3Client, params)
	var i int
	for p.HasMorePages() {
		i++
		page, err := p.NextPage(context.TODO())
		if err != nil {
			s.logger.Warn("failed to get page %v, %v", i, err)
		}

		res := []*proto.DataObject{}
		for _, obj := range page.Contents {
			arn := fmt.Sprintf(`arn:aws:s3:::%s/%s`, bucket, *obj.Key)
			lastModified := ""
			if obj.LastModified != nil {
				lastModified = obj.LastModified.Format("2006-01-02 15:04:05")
			}

			res = append(res, &proto.DataObject{
				RemoteId:     arn,
				ResourceName: *obj.Key,
				Uri:          arn,
				Metadata:     fmt.Sprintf(`{"last_modified": "%s"}`, lastModified)})
		}
		_ = cb.Callback(&proto.SyncResponse{Response: res})
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	conn := &S3Connector{
		logger: logger,
	}
	var pluginMap = map[string]plugin.Plugin{
		"connector": &connector.ConnectorGRPCPlugin{Impl: conn},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
