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

func (d *Updater) ReplaceBinary() {

}
