package download

import (
	"testing"

	"github.com/ra341/glacier/pkg/logger"
)

func init() {
	logger.InitConsole("debug", true)
}

func TestDownload(t *testing.T) {
	//srv := New("http://localhost:6699", nil)
	//
	//err := srv.Download(1, "./tmp/download", false, false)
	//require.NoError(t, err)

	//err = srv.Download(12, "./download")
	//require.NoError(t, err)

}
