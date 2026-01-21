package download

import (
	"testing"

	"github.com/ra341/glacier/pkg/logger"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.InitDefault()
}

func TestDownload(t *testing.T) {
	srv := New("http://localhost:6699")

	err := srv.Download(1, "./tmp/download")
	require.NoError(t, err)

	//err = srv.Download(12, "./download")
	//require.NoError(t, err)

}
