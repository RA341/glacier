package types

import (
	v1 "github.com/ra341/glacier/generated/library/v1"
)

func (g *Download) ToProto() *v1.Download {
	return &v1.Download{
		Client:       g.Client,
		DownloadId:   g.DownloadId,
		State:        g.State.String(),
		Progress:     g.Progress,
		DownloadPath: g.DownloadPath,
		DownloadUrl:  g.DownloadUrl,
	}
}

func (g *Download) FromProto(rpcDownload *v1.Download) {
	g.Client = rpcDownload.Client
	g.DownloadPath = rpcDownload.DownloadPath
	g.DownloadUrl = rpcDownload.DownloadUrl
}
