package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pidanou/c1-core/pkg/connector"
	"github.com/pidanou/c1-core/pkg/connector/proto"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GDriveConnector struct {
	logger hclog.Logger
}

type Options struct {
	CredentialsFile string `json:"credentials_file"`
}

func (s *GDriveConnector) Sync(options string, cb connector.CallbackInterface) (string, error) {

	var opts Options

	err := json.Unmarshal([]byte(options), &opts)
	if err != nil {
		s.logger.Error("Failed to unmarshal options", "error", err)
	}

	// Create GDrive service client
	ctx := context.Background()
	var svc *drive.Service
	svc, err = drive.NewService(ctx, option.WithCredentialsFile(opts.CredentialsFile))

	var pageToken string

	for {
		req := svc.Files.List().PageSize(1000).Fields("nextPageToken,files(id, name, webViewLink, modifiedTime)")

		if pageToken != "" {
			req.PageToken(pageToken)
		}

		resp, err := req.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
			return fmt.Sprintf("{\"message\":\"%v\"}", err), err
		}

		res := []*proto.DataObject{}
		for _, obj := range resp.Files {
			lastModified := ""
			if obj.ModifiedTime != "" {
				parsedTime, err := time.Parse(time.RFC3339, obj.ModifiedTime)
				if err != nil {
					lastModified = ""
				}
				lastModified = parsedTime.Format("2006-01-02 15:04:05")
			}

			res = append(res, &proto.DataObject{
				RemoteId:     obj.Id,
				ResourceName: obj.Name,
				Uri:          obj.WebViewLink,
				Metadata:     fmt.Sprintf(`{"last_modified": "%s"}`, lastModified)})
		}
		// Ignore proto.Empty, error response
		_ = cb.Callback(&proto.SyncResponse{Response: res})
		pageToken = resp.NextPageToken
		if resp.NextPageToken == "" {
			break
		}
	}

	return fmt.Sprintf("{\"message\":\"success\"}"), nil
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

	conn := &GDriveConnector{
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
