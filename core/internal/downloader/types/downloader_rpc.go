package types

import (
	"strings"

	v1 "github.com/ra341/glacier/generated/library/v1"
	"github.com/rs/zerolog/log"
)

func (g *Download) ToProto() *v1.Download {
	return &v1.Download{
		Client:       g.Client.String(),
		DownloadId:   g.DownloadId,
		State:        strings.TrimPrefix(g.State.String(), "Download"),
		Progress:     g.Progress,
		DownloadPath: g.DownloadPath,
		DownloadUrl:  g.DownloadUrl,
	}
}

func (g *Download) FromProto(rpcDownload *v1.Download) {
	clientType, err := ClientTypeString(rpcDownload.Client)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse client type")
		clientType = ClientUnknown
	}

	downloadType, err := DownloadStateString(rpcDownload.State)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse download type")
		downloadType = Unknown
	}

	g.Client = clientType
	g.DownloadId = rpcDownload.DownloadId
	g.State = downloadType
	g.Progress = rpcDownload.Progress
	g.DownloadPath = rpcDownload.DownloadPath
	g.DownloadUrl = rpcDownload.DownloadUrl
}
