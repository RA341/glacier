package download

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/ra341/glacier/pkg/logger"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.InitForTest()
}

func TestDownload(t *testing.T) {
	conf := Config{
		httpCli: new(http.Client),
		getBase: func() string {
			return "http://localhost:6699/api/server/protected"
		},
		SpeedThrottleInMB:       0,
		MaxConcurrentGames:      1,
		MaxConcurrentFiles:      10,
		MaxConcurrentFileChunks: 50,
		ChunkSizeInMB:           128,
	}

	editStatus := func(ctx context.Context, id int, down *LocalDownload) error {
		fmt.Println("editStatus", down)
		return nil
	}
	onDone := func(id int) {

	}
	download, err := NewDownload(
		&conf,
		editStatus,
		onDone,
		1,
		"./test/download",
		true,
	)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Go(download.Start)

	wg.Wait()
}
