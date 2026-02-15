package updater

import (
	"context"

	"github.com/google/go-github/v82/github"
)

const owner = "RA341"
const repo = "glacier"

type Updater struct {
	cli *github.Client
}

func New() *Updater {
	return &Updater{
		cli: github.NewClient(nil),
	}
}

func (d *Updater) ReplaceBinary(pathToNewBinary string) error {
	// oldExe running
	// oldExe execs pathToNewBinary --update --path=pathToOldBin
	// oldExe stops

	// newBin running in update mode
	// newBin renames pathToOldBin.old
	// newBin transfers its files using os.Executable and os read to pathToOldBin
	// newBin calls pathToOldbin --cleanUpdate=pathToUpdateFolder
	// newBin exits
	// replacedOldbin is now running

	return nil
}

func (d *Updater) CheckForUpdate() {

}

// FetchReleases fetches all releases for a specific repository
func (d *Updater) FetchReleases(ctx context.Context) ([]*github.RepositoryRelease, error) {
	releases, _, err := d.cli.Repositories.ListReleases(
		ctx,
		owner,
		repo,

		&github.ListOptions{PerPage: 10},
	)
	return releases, err
}

func (d *Updater) DownloadUpdate(downloadUrl string) {

}
