package util

import (
	"context"
	"go-micro.dev/v4/metadata"
)

type IpInfo struct {
	LocalIp  string
	RemoteIp string
}

// GetIp return a IpInfo from content.Content
func GetIp(ctx context.Context) (ipInfo IpInfo) {
	md, _ := metadata.FromContext(ctx)
	if _, ok := md["Local"]; ok {
		ipInfo.LocalIp = md["Local"]
	}
	if _, ok := md["Local"]; ok {
		ipInfo.RemoteIp = md["Remote"]
	}
	return
}
