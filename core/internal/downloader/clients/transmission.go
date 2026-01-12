package clients

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/hekmon/transmissionrpc/v3"
)

type TransmissionConfig struct {
	ClientType string `yaml:"clientType" config:"flag=ct,env=CLIENT_TYPE,default=,usage=Client: qbit|transmission|deluge"`
	Protocol   string `yaml:"protocol" config:"flag=torProto,env=TORRENT_PROTOCOL,default=http,usage=Protocol for torrent client connection http|https"`
	Host       string `yaml:"host" config:"flag=host,env=TORRENT_HOST,default=localhost:8080,usage=Host e.g. localhost:8080|qbit.somedomain.com"`
	User       string `yaml:"user" config:"flag=torUser,env=TORRENT_USER,default=,usage=Username for torrent client authentication"`
	Password   string `yaml:"password" config:"flag=torPass,env=TORRENT_PASSWORD,default=,usage=Password for torrent client authentication,hide=true"`
}
type TransmissionClient struct {
	Client *transmissionrpc.Client
}

func NewTransmissionClient(client *TransmissionConfig) (Downloader, error) {
	clientStr := ""
	if client.User != "" && client.Password != "" {
		clientStr = fmt.Sprintf("%s://%s:%s@%s/transmission/rpc", client.Protocol, client.User, client.Password, client.Host)
	} else {
		clientStr = fmt.Sprintf("%s://%s/transmission/rpc", client.Protocol, client.Host)
	}

	endpoint, err := url.Parse(clientStr)
	if err != nil {
		return nil, err
	}

	tbt, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		panic(err)
	}

	transmission := &TransmissionClient{Client: tbt}

	_, _, err = transmission.test()
	if err != nil {
		log.Error().Err(err).Msg("Transmission client check failed")
		return nil, fmt.Errorf("transmission client health check failed: %s", err)
	}

	return transmission, nil
}

func (tm *TransmissionClient) test() (string, string, error) {
	ok, serverVersion, serverMinimumVersion, err := tm.Client.RPCVersion(context.Background())
	if err != nil {
		return "", "", err
	}

	if !ok {
		return "", "", fmt.Errorf("remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion)
	}

	return "transmission", fmt.Sprintf("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)\n",
		serverVersion, transmissionrpc.RPCVersion), nil
}

func (tm *TransmissionClient) Download(ctx context.Context, url string, downloadPath string) (downloadID string, err error) {
	if strings.HasPrefix(url, "magnet:?") {
		log.Debug().Str("url", url).Msg("Downloading with magnet magnet link")
		return tm.downloadMagnet(ctx, url, downloadPath, "")
	}

	log.Debug().Str("url", url).Msg("Downloading with torrent file")
	return tm.downloadFile(url, downloadID, "")
}

func (tm *TransmissionClient) downloadFile(torrent, downloadPath, _ string) (string, error) {
	torrentResult, err := tm.Client.TorrentAddFileDownloadDir(context.Background(), torrent, downloadPath)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(*torrentResult.ID, 10), nil
}

func (tm *TransmissionClient) downloadMagnet(ctx context.Context, magnet, downloadPath, category string) (string, error) {
	torrentResult, err := tm.Client.TorrentAdd(ctx, transmissionrpc.TorrentAddPayload{
		Filename:    &magnet,
		DownloadDir: &downloadPath,
		Labels:      []string{category},
	})
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(*torrentResult.ID, 10), nil
}

func (tm *TransmissionClient) Cancel(ctx context.Context, downloadId string, removeDownloaded string) error {
	return fmt.Errorf(" this function is not implemented you idiot")
}

func (tm *TransmissionClient) Progress(ctx context.Context, downloadId string) (state DownloadState, progressString string, err error) {
	var intIds = make([]int64, 1)
	finalId, err := strconv.Atoi(downloadId)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to convert torrent id: %s to int64", downloadId)
		return Error, "", err
	}
	intIds = append(intIds, int64(finalId))

	infos, err := tm.Client.TorrentGetAllFor(context.Background(), intIds)
	if err != nil {
		return Error, "", err
	}
	if len(infos) == 0 {
		return Error, "", fmt.Errorf("no torrent not found")
	}

	for _, info := range infos {
		state = Downloading
		if info.Status.String() == "seeding" {
			state = Complete
		} else if info.Status.String() == "error" ||
			info.Status.String() == "stopped" {
			state = Error
		}

		progressString = fmt.Sprintf(
			"%s\nprogress:%v",
			info.Status.String(),
			info.PercentDone,
		)

		break
	}

	return state, progressString, nil
}
